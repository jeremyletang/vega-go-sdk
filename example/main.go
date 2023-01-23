package main

import (
	"context"
	"log"

	"github.com/jeremyletang/vega-go-sdk/wallet"

	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
)

const (
	token      = "SBmcqwG5jTJPqggNVSHddcWT07dHDV5EvXlAiQa1nYFXtD0bura87LIa9yPCDfAu"
	pubkey     = "ad2e531441c2e8a43e85423db399a4acc8f9a8a2376304a4c377d0da8eb31e80"
	walletAddr = "http://127.0.0.1:1789"
)

func main() {
	client, err := wallet.NewClient(walletAddr, token)
	if err != nil {
		log.Fatalf("unable to startup client: %v", err)
	}

	err = client.SendTransaction(context.Background(), pubkey, &walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &commandspb.VoteSubmission{
				ProposalId: "90e71c52b2f40db78efc24abe4217382993868cd24e45b3dd17147be4afaf884",
				Value:      vega.Vote_VALUE_NO,
			},
		},
	})

	if err != nil {
		log.Fatalf("unable to send transaction: %v", err)
	}
}
