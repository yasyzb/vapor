package bench

import (
	"encoding/json"

	"github.com/vapor/errors"
)

type Client struct {
	IP string
}

type response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrDetail string          `json:"error_detail"`
}

type submitTxReq struct {
	Tx interface{} `json:"raw_transaction"`
}

type submitTxResp struct {
	TxID string `json:"tx_id"`
}

func (c *Client) request(url string, payload []byte, respData interface{}) error {
	resp := &response{}
	if err := post(c.IP+url, payload, resp); err != nil {
		return err
	}

	if resp.Status != "success" {
		return errors.New(resp.ErrDetail)
	}

	return json.Unmarshal(resp.Data, respData)
}

func (c *Client) SubmitTx(tx interface{}) (string, error) {
	url := "/submit-transaction"
	payload, err := json.Marshal(submitTxReq{Tx: tx})
	if err != nil {
		return "", errors.Wrap(err, "json marshal")
	}

	res := &submitTxResp{}
	return res.TxID, c.request(url, payload, res)
}
