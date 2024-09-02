/*
Copyright © 2024 Nic Gibson <nic.gibson@redis.com>
*/
package cmd

import (
	"os"

	"github.com/goslogan/clusterinfo"
	"github.com/spf13/cobra"
)

var clusterInfo *clusterinfo.ClusterInfo

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rlatool",
	Short: "Tool for working with Redis Inc's rladmin tools's status output",
	Long: `rlatool parses, merges and outputs data from Redis Inc's rladmin tools.
Specifically, it takes the output of *rladmin status* and generates CSV and JSON
files from it. These can be used for further analysis.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		infile, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		if infile != "" {
			reader, err := os.Open(infile)
			if err != nil {
				return err
			} else {
				cmd.SetIn(reader)
			}
		}

		outfile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		if outfile != "" {
			writer, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			} else {
				cmd.SetOut(writer)
			}
		}

		clusterInfo, err = clusterinfo.NewClusterInfo(cmd.InOrStdin())
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("input", "i", "", "The rladmin status output to be parsed; if not provided STDIN is assumed")
	rootCmd.PersistentFlags().StringP("output", "o", "", "The file to which output should be written; if not provided STDOUT is assumed")
	rootCmd.MarkFlagFilename("input")
	rootCmd.MarkFlagFilename("output")
	rootCmd.InOrStdin()
}
