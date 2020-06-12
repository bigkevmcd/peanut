package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/peanut/pkg/parser"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func logIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peanut",
		Short: "Just a Go Kubernetes resource analyzer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := parser.Parse(viper.GetString("kustomization-path"))
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			defer w.Flush()

			for _, app := range cfg.Apps {
				fmt.Fprintf(w, "application: %s\n", app.Name)
				fmt.Fprintln(w, "name\tnamespace\treplicas\timages\t")
				for _, svc := range app.Services {
					images := strings.Join(svc.Images, ",")
					fmt.Fprintf(w, "%s\t%s\t%d\t%s\t\n", svc.Name, svc.Namespace, svc.Replicas, images)
				}
			}

			return nil
		},
	}

	cmd.Flags().String(
		"kustomization-path",
		"./kustomization.yaml",
		"path to read kustomization from",
	)
	logIfError(viper.BindPFlag("kustomization-path", cmd.Flags().Lookup("kustomization-path")))
	logIfError(cmd.MarkFlagRequired("kustomization-path"))

	cmd.AddCommand(makeHTTPCmd())
	return cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute is the main entry point into this component.
func Execute() {
	if err := makeRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
