package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CustomHelpTemplate = `
Gen data to databases with configurable concurrency and batch size.

Usage:
  gendata [command] [flags]

Available Commands:
  mysql        Generate data to MySQL database (alias: sql)
  postgres     Generate data to PostgreSQL database (alias: pg, postgresql)
  clickhouse   Generate data to ClickHouse database (alias: ck)
  completion   Generate the autocompletion script for the specified shell
  help         Help about any command

Global Flags:
  --concurrency   int      Number of concurrent workers to use for data generation (default: 1)
  --batchsize     int      Number of records to process in each batch (default: 1000)
  --repeatcount   int      Number of times to repeat the batch (default: 10)
                           Total records generated = batchsize * repeatcount * concurrency (default: 10000)
  --debug         bool     Enable debug log output

Database Connection Flags (required for mysql, postgres, clickhouse):
  --host             string          Database host
  --port             int             Database port
  --user             string          Database username
  --password         string          Database password
  --dbname           string          Database name
  --table            string          Database table (default: gendata_table)
  									 The actual table and its structure are fixed;
									 This option changes the default table name used for insertions
  --additionalargs   stringToString  Additional connection arguments (key=value format)
									 Test: --additionalargs "sslmode=disable,charset=utf8mb4"
Note:
The actual table and its structure are fixed; the --table option merely changes the default table name used for insertions.
`

// CustomHelpFunc 自定义帮助输出
func CustomHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", CustomHelpTemplate)
}
