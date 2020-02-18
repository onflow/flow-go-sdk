# sconfig
Configure your Go applications with a single struct

## Example

```go
import "github.com/psiemens/sconfig"

type Config struct {
    Environment string        `default:"LOCAL" flag:"env,e" info:"application environment"`
    Port        int           `default:"80" flag:"port,p" info:"port to start the server on"`
    Timeout     time.Duration `default:"1s"`
}

var conf Config

var cmd = &cobra.Command{
	Use: "hello-world",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Hello world!")
	},
}

func init() {
	err := sconfig.New(&conf).
		FromEnvironment("APP").
		BindFlags(cmd.PersistentFlags()).
		Parse()
	if err != nil {
		log.Fatal(err)
	}
}
```

```bash
# call with env vars and/or command line flags
APP_TIMEOUT=5s hello-world --port 8080 -e PRODUCTION
```