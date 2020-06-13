package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/peanut/pkg/config"
	httpapi "github.com/bigkevmcd/peanut/pkg/http"
)

func makeHTTPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "provide a simple app API",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.ParseFile(viper.GetString("config"))
			if err != nil {
				return err
			}
			http.Handle("/", httpapi.NewRouter(cfg))
			listen := fmt.Sprintf(":%d", viper.GetInt("port"))
			log.Printf("listening %s\n", listen)
			return http.ListenAndServe(listen, nil)
		},
	}

	cmd.Flags().Int(
		"port",
		8080,
		"port to serve requests on",
	)
	logIfError(viper.BindPFlag("port", cmd.Flags().Lookup("port")))

	cmd.Flags().String(
		"config",
		"",
		"file to parse configuration from",
	)
	logIfError(viper.BindPFlag("config", cmd.Flags().Lookup("config")))
	logIfError(cmd.MarkFlagRequired("config"))
	return cmd
}
