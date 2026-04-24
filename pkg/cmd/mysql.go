package cmd

import (
	"fmt"

	"gendata/pkg/action"

	"github.com/spf13/cobra"
)

func newMysqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mysql",
		Aliases: []string{"sql"},
		Short:   "Generate data to MySQL database",
		Long:    "Generate test data and write it to a MySQL database with configurable concurrency and batch size.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Run("mysql"); err != nil {
				return fmt.Errorf("failed to generate MySQL data: %w", err)
			}
			return nil
		},
	}
	// 设置自定义帮助函数
	cmd.SetHelpFunc(CustomHelpFunc)
	addSqlDatabaseFlags(cmd)
	return cmd
}
