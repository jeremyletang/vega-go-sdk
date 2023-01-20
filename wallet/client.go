package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	clientSendTransaction = "client.send_transaction"
)

var (
	healthEndpoint = func(baseUrl string) string {
		return fmt.Sprintf("%v%v", strings.TrimSuffix(baseUrl, "/"), "/api/v2/health")
	}
	requestEndpoint = func(baseUrl string) string {
		return fmt.Sprintf("%v%v", strings.TrimSuffix(baseUrl, "/"), "/api/v2/requests")
	}

	defaultTimeout = 5 * time.Second
)

type Client struct {
	walletAddr string
	token      string

	c *http.Client
}

func NewClient(walletAddr, token string) (*Client, error) {
	c := &http.Client{Timeout: defaultTimeout}
	req, err := http.NewRequest("GET", healthEndpoint(walletAddr), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from health endpoint: %v")
	}

	return &Client{
		walletAddr: walletAddr,
		token:      token,
		c:          c,
	}, nil
}

func (c *Client) SendTransaction(ctx context.Context, pubkey string, tx proto.Message) error {
	payload := newSendTransactionRequest(clientSendTransaction, pubkey, tx)
	body, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestEndpoint(c.walletAddr), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", c.token)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status from health endpoint: %v")
	}

	return nil
}
