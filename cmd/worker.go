/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/darcygail/ether-explorer/internal/worker"
	"github.com/darcygail/ether-explorer/schema"
	"github.com/spf13/cobra"
)

var workerConfig schema.WorkConfig

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		worker.Run(workerConfig)
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.PersistentFlags().BoolVar(&workerConfig.Full, "full", false, "parse all blocks")
	workerCmd.PersistentFlags().StringVar(&workerConfig.RpcUrl, "rpcUrl", "", "rpc url")
	workerCmd.PersistentFlags().StringVar(&workerConfig.MongoUri, "mongoUri", "", "mongo uri")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
