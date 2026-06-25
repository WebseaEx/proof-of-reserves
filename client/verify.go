package client

import (
	"websea-zkmerkle-proof/service/verify_service"

	"github.com/spf13/cobra"
)

func VerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "verify proof",
	}

	cmd.AddCommand(
		VerifyCexCommand(),
		VerifyUserCommand(),
	)
	return cmd
}

func VerifyCexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cex",
		Short: "verify cex proof",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify_service.CexVerify()
		},
	}
	return cmd
}

func VerifyUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "verify user proof",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify_service.UserVerify()
		},
	}
	return cmd
}
