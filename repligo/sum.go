package cmd

import (
	"github.com/spf13/cobra"
)

// sumCmd represents the sum command
var sumCmd = &cobra.Command{
	Use:   "sum",
	Short: "Compute BLAKE2 sum for input stream or file",
	Long:  "Compute BLAKE2 sum for input stream or file",
	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Compute BLAKE2 checksum for input
	},
}

func init() {
	//RootCmd.AddCommand(sumCmd)
}
