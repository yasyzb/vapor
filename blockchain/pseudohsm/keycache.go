package pseudohsm

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	vcrypto "github.com/vapor/crypto"
	edchainkd "github.com/vapor/crypto/ed25519/chainkd"
)

// Minimum amount of time between cache reloads. This limit applies if the platform does
// not support change notifications. It also applies if the keystore directory does not
// exist yet, the code will attempt to create a watcher at most this often.
const minReloadInterval = 2 * time.Second

type keysByFile []XPub

func (s keysByFile) Len() int           { return len(s) }
func (s keysByFile) Less(i, j int) bool { return s[i].File < s[j].File }
func (s keysByFile) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// AmbiguousKeyError is returned when attempting to unlock
// an XPub for which more than one file exists.
type AmbiguousKeyError struct {
	Pubkey  string
	Matches []XPub
}

func (err *AmbiguousKeyError) Error() string {
	files := ""
	for i, a := range err.Matches {
		files += a.File
		if i < len(err.Matches)-1 {
			files += ", "
		}
	}
	return fmt.Sprintf("multiple keys match keys (%s)", files)
}

// keyCache is a live index of all keys in the keystore.
type keyCache struct {
	keydir   string
	watcher  *watcher
	mu       sync.Mutex
	all      keysByFile
	byPubs   map[vcrypto.XPubKeyer][]XPub
	throttle *time.Timer
}

func newKeyCache(keydir string) *keyCache {
	kc := &keyCache{
		keydir: keydir,
		byPubs: make(map[vcrypto.XPubKeyer][]XPub),
	}
	kc.watcher = newWatcher(kc)
	return kc
}

func (kc *keyCache) hasKey(xpub vcrypto.XPubKeyer) bool {
	kc.maybeReload()
	kc.mu.Lock()
	defer kc.mu.Unlock()
	fmt.Println("hasKey xpub:", xpub)
	fmt.Println("hasKey xpub type:", reflect.TypeOf(xpub))
	fmt.Println("hasKey xpub:", reflect.ValueOf(xpub))
	fmt.Println("hasKey len(kc.byPubs[xpub]):", len(kc.byPubs[xpub]))
	fmt.Println("hasKey kc:", kc)
	fmt.Println("hasKey done...")
	fmt.Println("hasKey return:", len(kc.byPubs[xpub]) > 0)
	return len(kc.byPubs[xpub]) > 0
}

func (kc *keyCache) hasAlias(alias string) bool {
	fmt.Println("hasAlias alias:", alias)
	xpubs := kc.keys()
	for _, xpub := range xpubs {
		if xpub.Alias == alias {
			return true
		}
	}
	return false
}

func (kc *keyCache) add(newKey XPub) {
	fmt.Println("add start...")
	fmt.Println("add newKey:", newKey)
	fmt.Println("add newKey type:", reflect.TypeOf(newKey).String())
	fmt.Println("add newKey.XPub type:", reflect.TypeOf(newKey.XPub).String())
	kc.mu.Lock()
	defer kc.mu.Unlock()

	// switch xpub := newKey.XPub.(type){
	// 	if reflect.TypeOf(xpub).
	// }

	i := sort.Search(len(kc.all), func(i int) bool { return kc.all[i].File >= newKey.File })
	if i < len(kc.all) && kc.all[i] == newKey {
		return
	}
	fmt.Println("add i:", i)
	fmt.Println("add kc.all 1:", kc.all)
	// newKey is not in the cache.
	kc.all = append(kc.all, XPub{})
	fmt.Println("add kc.all 2:", kc.all)
	copy(kc.all[i+1:], kc.all[i:])
	fmt.Println("add kc.all 3:", kc.all)
	kc.all[i] = newKey
	fmt.Println("add newKey:", newKey.XPub)
	fmt.Println("add newKey -4 type:", reflect.TypeOf(newKey.XPub).String())
	fmt.Println("add kc.all 4:", kc.all)
	kc.byPubs[newKey.XPub] = append(kc.byPubs[newKey.XPub], newKey)
	fmt.Println("add kc.byPubs:", kc.byPubs)
	fmt.Println("add kc.byPubs:", kc.byPubs[newKey.XPub])
	fmt.Println("add kc.byPubs type:", reflect.TypeOf(kc.byPubs).String())
	fmt.Println("add kc:", kc)
}

func (kc *keyCache) keys() []XPub {
	fmt.Println("keys start...")
	fmt.Println("keys kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("keys i:", i)
		fmt.Println("keys i type:", reflect.TypeOf(i))
		fmt.Println("keys v:", v)
		fmt.Println("keys v type:", reflect.TypeOf(v))
	}
	kc.maybeReload()
	kc.mu.Lock()
	defer kc.mu.Unlock()
	cpy := make([]XPub, len(kc.all))
	copy(cpy, kc.all)
	fmt.Println("keys end...")
	fmt.Println("keys kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("keys i:", i)
		fmt.Println("keys i type:", reflect.TypeOf(i))
		fmt.Println("keys v:", v)
		fmt.Println("keys v type:", reflect.TypeOf(v))
	}

	return cpy
}

func (kc *keyCache) maybeReload() {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if kc.watcher.running {
		return // A watcher is running and will keep the cache up-to-date.
	}

	if kc.throttle == nil {
		kc.throttle = time.NewTimer(0)
	} else {
		select {
		case <-kc.throttle.C:
		default:
			return // The cache was reloaded recently.
		}
	}
	kc.watcher.start()
	kc.reload()
	kc.throttle.Reset(minReloadInterval)
}

// find returns the cached keys for alias if there is a unique match.
// The exact matching rules are explained by the documentation of Account.
// Callers must hold ac.mu.
func (kc *keyCache) find(xpub XPub) (XPub, error) {
	fmt.Println("find start...")
	// Limit search to xpub candidates if possible.
	matches := kc.all
	fmt.Println("matches := kc.all", matches)
	// if (xpub.XPub != vcrypto.XPubKeyer{}) {
	// matches = kc.byPubs[xpub.XPub]
	// }
	// if _, ok := xpub.XPub.(edchainkd.XPub); ok {
	fmt.Println("find xpub.XPub:", xpub.XPub)
	fmt.Println("find xpub.XPub type:", reflect.TypeOf(xpub.XPub))
	matches = kc.byPubs[xpub.XPub]
	fmt.Println("find matches:", matches)
	// }

	for i, v := range matches {
		fmt.Println("find i:", i)
		fmt.Println("find i type:", reflect.TypeOf(i))
		fmt.Println("find v:", v)
		fmt.Println("find v type:", reflect.TypeOf(v))
	}

	// for i, c := range kc.byPubs {
	// 	if reflect.TypeOf(i).String() == "string" {
	// 		if xpb, err := edchainkd.NewXPub(reflect.ValueOf(i).String()); err != nil {
	// 			panic(err)
	// 		} else {
	// 			kc.byPubs[*xpb] = c
	// 			delete(kc.byPubs, i)
	// 		}
	// 	}
	// 	fmt.Println("LoadChainKDKey i:", i)
	// 	fmt.Println("LoadChainKDKey i type:", reflect.TypeOf(i))
	// 	fmt.Println("LoadChainKDKey c:", c)
	// 	fmt.Println("LoadChainKDKey c type:", reflect.TypeOf(c[0]))
	// }
	fmt.Println("xpub.XPub:", xpub.XPub)
	fmt.Println("xpub.XPub type:", reflect.TypeOf(xpub.XPub))
	fmt.Println("keyCache:", kc)
	fmt.Println("keyCache byPubs:", kc.byPubs)
	fmt.Println("keyCache byPubs xpub.XPub:", kc.byPubs[xpub.XPub])
	// switch xpb := xpub.XPub.(type) {
	// case edchainkd.XPub:
	// 	fmt.Println("ed25519 xpb")
	// 	matches = kc.byPubs[xpb]
	// 	fmt.Println("ed25519 mathches:", matches)
	// }
	// matches = kc.byPubs[xpub.XPub]
	fmt.Println("mathes:", matches)
	fmt.Println("find xpub.File:", xpub.File)
	if xpub.File != "" {
		// If only the basename is specified, complete the path.
		if !strings.ContainsRune(xpub.File, filepath.Separator) {
			xpub.File = filepath.Join(kc.keydir, xpub.File)
		}
		for i := range matches {
			if matches[i].File == xpub.File {
				return matches[i], nil
			}
		}
		if _, ok := xpub.XPub.(edchainkd.XPub); ok {
			return XPub{}, ErrLoadKey
		}
		// or other, e.g. sm2

		return XPub{}, ErrLoadKey
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return XPub{}, ErrLoadKey
	default:
		err := &AmbiguousKeyError{Pubkey: xpubToString(xpub.XPub), Matches: make([]XPub, len(matches))}
		copy(err.Matches, matches)
		return XPub{}, err
	}
}

func xpubToString(xpub vcrypto.XPubKeyer) (str string) {
	switch xpbk := xpub.(type) {
	case edchainkd.XPub:
		return hex.EncodeToString(xpbk[:])
	}

	return
}

// reload caches addresses of existing key.
// Callers must hold ac.mu.
func (kc *keyCache) reload() {
	fmt.Println("reload before scan.......")
	fmt.Println("reload kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("reload i:", i)
		fmt.Println("reload i type:", reflect.TypeOf(i))
		fmt.Println("reload v:", v)
		fmt.Println("reload v type:", reflect.TypeOf(v))
	}

	keys, err := kc.scan()
	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "load keys error": err}).Error("can't load keys")
	}
	kc.all = keys
	fmt.Println("reload before sort.......")
	fmt.Println("reload kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("reload i:", i)
		fmt.Println("reload i type:", reflect.TypeOf(i))
		fmt.Println("reload v:", v)
		fmt.Println("reload v type:", reflect.TypeOf(v))
	}
	sort.Sort(kc.all)
	for k := range kc.byPubs {
		delete(kc.byPubs, k)
	}
	fmt.Println("reload after sort.......")
	fmt.Println("reload kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("reload i:", i)
		fmt.Println("reload i type:", reflect.TypeOf(i))
		fmt.Println("reload v:", v)
		fmt.Println("reload v type:", reflect.TypeOf(v))
	}
	for _, k := range keys {
		// xpub := edchainkd.XPub{}
		switch kt := k.XPub.(type) {
		case string:
			xpb, err := edchainkd.NewXPub(kt)
			if err != nil {
				panic(err)
			}
			k.XPub = *xpb
		}
		kc.byPubs[k.XPub] = append(kc.byPubs[k.XPub], k)
	}
	fmt.Println("reload end.......")
	fmt.Println("reload kc:", kc)
	for i, v := range kc.byPubs {
		fmt.Println("reload i:", i)
		fmt.Println("reload i type:", reflect.TypeOf(i))
		fmt.Println("reload v:", v)
		fmt.Println("reload v type:", reflect.TypeOf(v))
	}
	log.WithFields(log.Fields{"module": logModule, "cache has keys:": len(kc.all)}).Debug("reloaded keys")
}

func (kc *keyCache) scan() ([]XPub, error) {
	files, err := ioutil.ReadDir(kc.keydir)
	if err != nil {
		return nil, err
	}
	var (
		buf     = new(bufio.Reader)
		keys    []XPub
		keyJSON struct {
			Alias string            `json:"alias"`
			XPub  vcrypto.XPubKeyer `json:"xpub"`
		}
	)
	for _, fi := range files {
		path := filepath.Join(kc.keydir, fi.Name())
		if skipKeyFile(fi) {
			//log.Printf("ignoring file %v", path)
			//fmt.Printf("ignoring file %v", path)
			continue
		}
		fd, err := os.Open(path)
		if err != nil {
			//log.Printf(err)
			fmt.Printf("err")
			continue
		}
		buf.Reset(fd)
		// Parse the address.
		keyJSON.Alias = ""
		err = json.NewDecoder(buf).Decode(&keyJSON)
		switch {
		case err != nil:
			log.WithFields(log.Fields{"module": logModule, "decode json err": err}).Errorf("can't decode key %s: %v", path, err)

		case (keyJSON.Alias == ""):
			log.WithFields(log.Fields{"module": logModule, "can't decode key, key path:": path}).Warn("missing or void alias")
		default:
			keys = append(keys, XPub{XPub: keyJSON.XPub, Alias: keyJSON.Alias, File: path})
		}
		fd.Close()
	}
	return keys, err
}

func (kc *keyCache) delete(removed XPub) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	kc.all = removeKey(kc.all, removed)
	if ba := removeKey(kc.byPubs[removed.XPub], removed); len(ba) == 0 {
		delete(kc.byPubs, removed.XPub)
	} else {
		kc.byPubs[removed.XPub] = ba
	}
}

func removeKey(slice []XPub, elem XPub) []XPub {
	for i := range slice {
		if slice[i] == elem {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func skipKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

func (kc *keyCache) close() {
	kc.mu.Lock()
	kc.watcher.close()
	if kc.throttle != nil {
		kc.throttle.Stop()
	}
	kc.mu.Unlock()
}
