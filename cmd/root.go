package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)
}

var RootCmd = &cobra.Command{
	Use:   "injun",
	Short: "injun: Search microservice for Eventackle",
	Long:  "injun: Search microservice for Eventackle backed by ElasticSearch",
	Version: "0.0.1",
}
