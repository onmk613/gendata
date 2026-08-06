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
  --concurrency      int       Number of concurrent write workers (default: 1)
  --batchsize        int       Number of records to generate per write batch (default: 1000)
                               Large batches are split automatically to stay within DB limits
  --repeatcount      int       Number of times to repeat the batch write (default: 10)
                               Used only in write mode when --duration = 0;
                               total records = batchsize * repeatcount * concurrency
  --mode             string    Run mode: write | read | mixed (default: write)
                               write  = pure insert throughput test
                               read   = point select by user_id (requires --duration)
                               mixed  = concurrent read + write (requires --duration)
  --readconcurrency   int      Number of concurrent read workers (default: same as --concurrency)
                               Read/write ratio is controlled by concurrency vs readconcurrency
  --duration         duration  Run for a fixed duration, e.g. 60s, 2m (default: 0 = disabled)
                               Required for read/mixed mode; overrides --repeatcount in write mode
  --debug            bool      Enable debug log output

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
Read mode never creates or alters the table; prepare data with write mode first.
`

// CustomHelpFunc 自定义帮助输出
func CustomHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", CustomHelpTemplate)
}
