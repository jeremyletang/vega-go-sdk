package wallet

import (
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type request struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      string          `json:"id,omitempty"`
}

type sendTransactionRequest struct {
	Pubkey      string          `json:"publicKey"`
	SendingMode string          `json:"sendingMode"`
	Transaction json.RawMessage `json:"transaction"`
}

func newSendTransactionRequest(method, pubkey string, msg proto.Message) *request {
	m := jsonpb.Marshaler{Indent: "  "}
	buf, _ := m.MarshalToString(msg)

	params := sendTransactionRequest{
		Pubkey:      pubkey,
		SendingMode: "TYPE_SYNC",
		Transaction: []byte(buf),
	}

	rawParams, _ := json.Marshal(params)

	return &request{
		Version: "2",
		Method:  method,
		Params:  rawParams,
	}
}
