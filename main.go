package main

import (
	"fmt"
	"os"
	"websea-zkmerkle-proof/client"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gproof",
		Short: "Command line interface for interacting with websea-zkmerkle-proof",
	}

	rootCmd.AddCommand(
		client.KeygenCommand(),
		client.WitnessCommand(),
		client.ProverCommand(),
		client.UserProofCommand(),
		client.VerifyCommand(),
		client.ToolCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
