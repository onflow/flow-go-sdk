package start

import (
	"fmt"
	"os"
	"time"

	"github.com/psiemens/sconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/cli"
	"github.com/dapperlabs/flow-go-sdk/cli/flow/initialize"
	"github.com/dapperlabs/flow-go-sdk/emulator/server"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

type Config struct {
	Port      int           `default:"3569" flag:"port,p" info:"port to run RPC server"`
	HTTPPort  int           `default:"8080" flag:"http-port" info:"port to run HTTP server"`
	Verbose   bool          `default:"false" flag:"verbose,v" info:"enable verbose logging"`
	BlockTime time.Duration `flag:"block-time,b" info:"time between sealed blocks"`
	RootKey   string        `flag:"root-key,k" info:"root account key"`
	Init      bool          `default:"false" flag:"init" info:"whether to initialize a new account profile"`
	GRPCDebug bool          `default:"false" flag:"grpc-debug" info:"enable gRPC server reflection for debugging with grpc_cli"`
	DBPath    string        `flag:"db" info:"where Flow chain data will be stored"`
}

var (
	log  *logrus.Logger
	conf Config
)

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Flow emulator server",
	Run: func(cmd *cobra.Command, args []string) {
		if conf.Init {
			pconf := initialize.InitProject()
			rootAcct := pconf.RootAccount()

			fmt.Printf("⚙️   Flow client initialized with root account:\n\n")
			fmt.Printf("👤  Address: 0x%x\n", rootAcct.Address.Bytes())
		}

		var rootKey flow.AccountPrivateKey

		if len(conf.RootKey) > 0 {
			rootKey = keys.MustDecodePrivateKeyHex(conf.RootKey)
		} else {
			rootAcct := cli.LoadConfig().RootAccount()
			rootKey = rootAcct.PrivateKey
		}

		if conf.Verbose {
			log.SetLevel(logrus.DebugLevel)
		}

		serverConf := &server.Config{
			GRPCPort:  conf.Port,
			GRPCDebug: conf.GRPCDebug,
			HTTPPort:  conf.HTTPPort,
			// TODO: allow headers to be parsed from environment
			HTTPHeaders:    nil,
			BlockTime:      conf.BlockTime,
			RootAccountKey: &rootKey,
			DBPath:         conf.DBPath,
		}

		emu := server.NewEmulatorServer(log, serverConf)
		emu.Start()
	},
}

func init() {
	initLogger()
	initConfig()
}

func initLogger() {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Out = os.Stdout
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
