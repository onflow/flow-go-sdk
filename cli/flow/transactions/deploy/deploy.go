package deploy

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/psiemens/sconfig"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/cli"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
	utils "github.com/dapperlabs/flow-go-sdk/utils/examples"
)

type Config struct {
	Signer string `default:"root" flag:"signer,s"`
	Host   string `default:"127.0.0.1:3569" flag:"host" info:"Flow Observation API host address"`
}

var conf Config

var Cmd = &cobra.Command{
	Use:   "deploy [path to contract]",
	Short: "Deploy a contract",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectConf := cli.LoadConfig()

		signer := projectConf.Accounts[conf.Signer]

		contractPath := args[0]

		contract, err := ioutil.ReadFile(contractPath)
		if err != nil {
			cli.Exitf(1, "Failed to load Cadence code from %s", contractPath)
		}

		deployScript, _ := templates.CreateAccount(nil, contract)

		tx := flow.Transaction{
			Script:       deployScript,
			Nonce:        rand.Uint64(),
			ComputeLimit: 10,
			PayerAccount: signer.Address,
		}

		sig, err := keys.SignTransaction(tx, signer.PrivateKey)
		if err != nil {
			cli.Exit(1, "Failed to sign transaction")
		}

		tx.AddSignature(signer.Address, sig)

		client, err := client.New(conf.Host)
		if err != nil {
			cli.Exit(1, "Failed to connect to emulator")
		}

		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			cli.Exitf(1, "Failed to send transaction: %v", err)
		}

		deployContractTxResp := utils.WaitForSeal(context.Background(), client, tx.Hash())

		var contractAddress flow.Address

		for _, event := range deployContractTxResp.Events {
			if event.Type == flow.EventAccountCreated {
				accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(event.Payload)
				utils.Handle(err)

				contractAddress = accountCreatedEvent.Address()
			}
		}

		fmt.Printf("Contract deployed to: 0x%s\n", contractAddress.Hex())
	},
}

func init() {
	initConfig()
}

func initConfig() {
	err := sconfig.New(&conf).
		FromEnvironment(cli.EnvPrefix).
		BindFlags(Cmd.PersistentFlags()).
		Parse()
	if err != nil {
		log.Fatal(err)
	}
}
