/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate a JSON conversion of rladmin status output",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := clusterInfo.JSON()

		if err == nil {
			fmt.Fprint(cmd.OutOrStdout(), output)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
