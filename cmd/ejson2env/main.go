package main

import (
	"fmt"
	"os"

	"github.com/Shopify/ejson2env/v2"
	"github.com/urfave/cli"
)

// version information. This will be overridden by the ldflags
var version = "dev"

// fail prints the error message to stderr, then ends execution.
func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Usage = "get environment variables from ejson files"
	app.Version = version
	app.Author = "Catherine Jones"
	app.Email = "catherine.jones@shopify.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "keydir, k",
			Value:  "/opt/ejson/keys",
			Usage:  "Directory containing EJSON keys",
			EnvVar: "EJSON_KEYDIR",
		},
		cli.BoolFlag{
			Name:  "key-from-stdin",
			Usage: "Read the private key from STDIN",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress export statement",
		},
		cli.BoolFlag{
			Name:  "trim-underscore",
			Usage: "Trim leading underscore from variable names",
		},
	}

	app.Action = func(c *cli.Context) {
		var filename string
		var userSuppliedPrivateKey string

		keydir := c.String("keydir")
		quiet := c.Bool("quiet")
		trim_underscore := c.Bool("trim-underscore")

		// select the ExportFunction to use
		exportFunc := ejson2env.ExportEnv
		if quiet {
			exportFunc = ejson2env.ExportQuiet
		}

		if trim_underscore {
			exportFunc = ejson2env.TrimLeadingUnderscoreExportWrapper(exportFunc)
		}

		if c.Bool("key-from-stdin") {
			var err error
			userSuppliedPrivateKey, err = readKey(os.Stdin)
			if err != nil {
				fail(fmt.Errorf("failed to read from stdin: %s", err))
			}
		}

		if 1 <= len(c.Args()) {
			filename = c.Args().Get(0)
		}

		if "" == filename {
			fail(fmt.Errorf("no secrets.ejson filename passed"))
		}

		if err := ejson2env.ReadAndExportEnv(filename, keydir, userSuppliedPrivateKey, exportFunc); nil != err {
			fail(err)
		}
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Unexpected failure:", err)
		os.Exit(1)
	}

}
