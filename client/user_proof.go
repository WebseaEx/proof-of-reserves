package client

import (
	"websea-zkmerkle-proof/service/userproof_service"

	"github.com/spf13/cobra"
)

func UserProofCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "userproof",
		Short: "used to generate and persist user merkle proof",
		Run: func(cmd *cobra.Command, args []string) {
			userproof_service.Handler()
		},
	}
	return cmd
}
