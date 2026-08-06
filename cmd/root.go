package cmd

import (
	"gendata/internal/action"

	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "gendata",
		Short:        "Gen data to databases",
		SilenceUsage: true,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			setupLogger()
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			action.CloseGendata()
		},
	}

	// 设置自定义帮助函数
	cmd.SetHelpFunc(CustomHelpFunc)

	flags := cmd.PersistentFlags()
	addRootFlags(flags)

	cmd.AddCommand(newMysqlCmd(), newPostgresCmd(), newClickhouseCmd())
	return cmd, nil
}
