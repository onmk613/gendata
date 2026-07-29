package cmd

import (
	"fmt"

	"gendata/internal/action"

	"github.com/spf13/cobra"
)

func newPostgresCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgres",
		Aliases: []string{"pg", "postgresql"},
		Short:   "Generate data to PostgreSQL database",
		Long:    "Generate test data and write it to a PostgreSQL database with configurable concurrency and batch size.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Run("postgres"); err != nil {
				return fmt.Errorf("failed to generate PostgreSQL data: %w", err)
			}
			return nil
		},
	}
	// 设置自定义帮助函数
	cmd.SetHelpFunc(CustomHelpFunc)
	addSqlDatabaseFlags(cmd)
	return cmd
}
