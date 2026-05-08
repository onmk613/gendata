package cmd

import (
	"github.com/spf13/cobra"

	"gendata/pkg/action"
)

func NewRootCmd(args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "gendata",
		Short:        "Gen data to databases",
		SilenceUsage: true,
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			action.CloseGendata()
		},
	}

	// 设置自定义帮助函数
	cmd.SetHelpFunc(CustomHelpFunc)

	flags := cmd.PersistentFlags()
	addRootFlags(flags)

	flags.ParseErrorsAllowlist.UnknownFlags = true
	flags.Parse(args)

	// 在flag解析之后立即设置日志
	setupLogger()

	cmd.AddCommand(newMysqlCmd(), newPostgresCmd(), newClickhouseCmd())
	return cmd, nil
}
