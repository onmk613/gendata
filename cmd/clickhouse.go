package cmd

import (
	"fmt"

	"gendata/internal/action"

	"github.com/spf13/cobra"
)

func newClickhouseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clickhouse",
		Aliases: []string{"ck"},
		Short:   "Generate data to ClickHouse database",
		Long:    "Generate test data and write it to a ClickHouse database with configurable concurrency and batch size.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Run("clickhouse"); err != nil {
				return fmt.Errorf("failed to generate ClickHouse data: %w", err)
			}
			return nil
		},
	}
	// 设置自定义帮助函数
	cmd.SetHelpFunc(CustomHelpFunc)
	addSqlDatabaseFlags(cmd)
	return cmd
}
