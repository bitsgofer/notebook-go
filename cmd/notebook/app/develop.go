package app

import (
	"fmt"

	"github.com/bitsgofer/notebook-go/internal/webserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "Run render on file changes and serve over HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverAddr := viper.GetString(cfgKeyDevServerAddr)
		dataDir := viper.GetString(cfgKeyDevServerDataDir)

		fmt.Printf("Listen to %s and serve pages from: %s\n", serverAddr, dataDir)
		srv, _ := webserver.NewDevelopmentServer(serverAddr, dataDir)

		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("=> Cannot run web server; err= %q\n", err)
			return nil
		}
		return nil
	},
}
