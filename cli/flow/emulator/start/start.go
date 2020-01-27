package start

import (
	"fmt"
	"os"
	"time"

	"github.com/psiemens/sconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/dapperlabs/flow-go-sdk/cli"
	"github.com/dapperlabs/flow-go-sdk/cli/flow/initialize"
	"github.com/dapperlabs/flow-go-sdk/emulator/server"
)

type Config struct {
	Port          int           `default:"3569" flag:"port,p" info:"port to run RPC server"`
	HTTPPort      int           `default:"8080" flag:"http_port" info:"port to run HTTP server"`
	Verbose       bool          `default:"false" flag:"verbose,v" info:"enable verbose logging"`
	BlockInterval time.Duration `default:"5s" flag:"interval,i" info:"time between minted blocks"`
	RootKey       string        `flag:"root-key" info:"root account key"`
	Init          bool          `default:"false" flag:"init" info:"whether to initialize a new account profile"`
	AutoMine      bool          `default:"true" flag:"automine" info:"enable instant transaction mining"`
	GRPCDebug     bool          `default:"false" flag:"grpc-debug" info:"enable gRPC server reflection for debugging with grpc_cli"`
	Persistent    bool          `default:"false" flag:"persistent" info:"enable persistent storage"`
	DBPath        string        `default:"./flowdb" flag:"db-path" info:"where Flow chain data will be stored"`
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
			var pconf *cli.Config
			if len(conf.RootKey) > 0 {
				prKey := cli.MustDecodeAccountPrivateKeyHex(conf.RootKey)
				pconf = initialize.InitProjectWithRootKey(prKey)
			} else {
				pconf = initialize.InitProject()
			}
			rootAcct := pconf.RootAccount()

			fmt.Printf("⚙️   Flow client initialized with root account:\n\n")
			fmt.Printf("👤  Address: 0x%x\n", rootAcct.Address.Bytes())
		}

		rootAcct := cli.LoadConfig().RootAccount()

		if conf.Verbose {
			log.SetLevel(logrus.DebugLevel)
		}

		serverConf := &server.Config{
			Port:           conf.Port,
			HTTPPort:       conf.HTTPPort,
			BlockInterval:  conf.BlockInterval,
			RootAccountKey: &rootAcct.PrivateKey,
			AutoMine:       conf.AutoMine,
			GRPCDebug:      conf.GRPCDebug,
			Persistent:     conf.Persistent,
			DBPath:         conf.DBPath,
		}

		server.StartServer(log, serverConf)
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
