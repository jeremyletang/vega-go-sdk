package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status from health endpoint: %v", resp.StatusCode)
	}

	return &Client{
		walletAddr: walletAddr,
		token:      token,
		c:          c,
	}, nil
}

func (c *Client) SendTransaction(ctx context.Context, pubkey string, tx *walletpb.SubmitTransactionRequest) error {
	payload := newSendTransactionRequest(clientSendTransaction, pubkey, tx)
	body, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestEndpoint(c.walletAddr), bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.addAuthToken(req)
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response: %v - %v", resp.StatusCode, string(content))
	}

	return nil
}

func (c *Client) addAuthToken(req *http.Request) {
	req.Header.Add("Origin", c.walletAddr)
	req.Header.Add("Authorization", fmt.Sprintf("VWT %v", c.token))
}
