package client

import (
	"websea-zkmerkle-proof/service/keygen_service"

	"github.com/spf13/cobra"
)

func KeygenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "generating zk related keys which are used to generate and verify zk proof",
		Run: func(cmd *cobra.Command, args []string) {
			keygen_service.Handler()
		},
	}
	return cmd
}
