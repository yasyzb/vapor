package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/golang/time/rate"
	"github.com/vapor/test/bench"
	"github.com/vapor/test/mock"
)

const (
	fedPkey = "e07a306fa454ab95b32f4d184effa87c0caf1cc2182bdbdd8e0207392787254d1967589f0d9dec13a388c0412002d2c267bdf3b920864e1ddc50581be5604ce1"
	fed     = "http://127.0.0.1:9889"
)

func main() {
	count := flag.Int("count", 0, "count")
	verbose := flag.Bool("v", true, "verbose")
	flag.Parse()

	c := &bench.Client{IP: fed}
	context := context.Background()
	limiter := rate.NewLimiter(10*1024, 10*1024)
	current := 0
	for {
		limiter.Wait(context)
		tx := mock.NewCrosschainTx(fedPkey)
		if tid, err := c.SubmitTx(tx); err != nil {
			fmt.Printf("%v: host: %v err: %v\n", current, c.IP, err)
		} else if *verbose {
			fmt.Printf("%v: host: %v tid: %v spent: %v\n", current, c.IP, tid, tx.Tx.MainchainOutputIDs[0])
		}
		current++
		if *count != 0 && current == *count {
			break
		}
		if current%1000 == 0 {
			fmt.Printf("%v\n", time.Now())
		}
	}
}
