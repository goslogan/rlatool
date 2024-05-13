package cmd

import (
	"fmt"

	"github.com/goslogan/rlatool/clusterinfo"
	"github.com/spf13/cobra"
)

func buildOutput(cmd *cobra.Command, input clusterinfo.Serializer) error {

	var output string

	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}

	if json {
		output, err = input.JSON()
	} else {
		output, err = input.CSV()
	}

	if err != nil {
		return err
	} else {
		fmt.Fprint(cmd.OutOrStdout(), output)
		return nil
	}

}
