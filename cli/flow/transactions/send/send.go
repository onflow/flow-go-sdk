package send

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/psiemens/sconfig"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/cli"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

type Config struct {
	Signer string `default:"root" flag:"signer,s"`
	Code   string `flag:"code,c"`
	Nonce  uint64 `flag:"nonce,n"`
	Host   string `default:"127.0.0.1:3569" flag:"host" info:"Flow Observation API host address"`
}

var conf Config

var Cmd = &cobra.Command{
	Use:   "send",
	Short: "Send a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		projectConf := cli.LoadConfig()

		signer := projectConf.Accounts[conf.Signer]

		var (
			code []byte
			err  error
		)

		if conf.Code != "" {
			code, err = ioutil.ReadFile(conf.Code)
			if err != nil {
				cli.Exitf(1, "Failed to load BPL code from %s", conf.Code)
			}
		}

		tx := flow.Transaction{
			Script:         code,
			Nonce:          conf.Nonce,
			ComputeLimit:   10,
			PayerAccount:   signer.Address,
			ScriptAccounts: []flow.Address{signer.Address},
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
