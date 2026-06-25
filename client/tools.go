package client

import (
	"websea-zkmerkle-proof/service/tool_service"

	"github.com/spf13/cobra"
)

func ToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tool",
		Short: "provide some common services",
	}

	cmd.AddCommand(
		ToolCleanKvrocks(),
		ToolCheckProverStatus(),
		ToolQueryCexAssets(),
		ToolCheckUserAssetsFile(),
		ToolCompareAccountID(),
	)
	return cmd
}

func ToolCleanKvrocks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean_kvrocks",
		Short: "remove only kvrocks data",
		Run: func(cmd *cobra.Command, args []string) {
			tool_service.CleanKvrocks()
		},
	}
	return cmd
}

func ToolCheckProverStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check_prover_status",
		Short: "check prover data status",
		Run: func(cmd *cobra.Command, args []string) {
			tool_service.CheckProverStatus()
		},
	}
	return cmd
}

func ToolQueryCexAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query_cex_assets",
		Short: "get cex assets info in json format",
		Run: func(cmd *cobra.Command, args []string) {
			tool_service.QueryCexAssets()
		},
	}
	return cmd
}

func ToolCheckUserAssetsFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check_user_assets_file",
		Short: "check user assets file csv",
		Run: func(cmd *cobra.Command, args []string) {
			tool_service.CheckUserFiles()
		},
	}
	return cmd
}

func ToolCompareAccountID() *cobra.Command {
	var filePath string
	var showAll bool
	var limit int

	cmd := &cobra.Command{
		Use:   "compare_account_id",
		Short: "compare csv uid with normalized account_id",
		Run: func(cmd *cobra.Command, args []string) {
			tool_service.CompareAccountID(filePath, showAll, limit)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "csv file path, defaults to UserDataFile in config")
	cmd.Flags().BoolVar(&showAll, "show-all", false, "print all rows instead of mismatches only")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum mismatch rows to print when --show-all is false")
	return cmd
}
