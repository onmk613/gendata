package cmd

import (
	"log/slog"
	"os"

	"gendata/internal/action"
	"gendata/internal/core"
	mydriver "gendata/internal/driver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var debugFlag bool

func addRootFlags(flags *pflag.FlagSet) {
	flags.IntVar(&action.WriteConf.Concurrency, "concurrency", 1, "Number of concurrent write workers")
	flags.IntVar(&action.WriteConf.BatchSize, "batchsize", 1000, "Number of records to write in each batch")
	flags.IntVar(&action.WriteConf.RepeatCount, "repeatcount", 10, "Number of times to repeat the batch write (ignored when --duration > 0 in write mode)")
	flags.StringVar(&action.WriteConf.Mode, "mode", "write", "Run mode: write | read | mixed")
	flags.IntVar(&action.WriteConf.ReadConcurrency, "readconcurrency", 0, "Number of concurrent read workers (0 = same as --concurrency)")
	flags.DurationVar(&action.WriteConf.Duration, "duration", 0, "Run for a fixed duration (e.g. 60s, 2m); required for read/mixed mode")
	flags.BoolVar(&debugFlag, "debug", false, "Print log")
}

func setupLogger() {
	if debugFlag {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.DiscardHandler))
	}
}

func addSqlDatabaseFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&mydriver.SqlConf.DBName, "dbname", "", "Database name")
	flags.StringVar(&mydriver.SqlConf.Host, "host", "", "Database host")
	flags.IntVar(&mydriver.SqlConf.Port, "port", 0, "Database port")
	flags.StringVar(&mydriver.SqlConf.User, "user", "", "Database user")
	flags.StringVar(&mydriver.SqlConf.Password, "password", "", "Database password")
	flags.StringVar(&core.TableName, "table", "gendata_table", "Database table")
	flags.StringToStringVar(&mydriver.SqlConf.AdditionalArgs, "additionalargs", nil, "Additional connection arguments (key=value format, multiple args separated by comma)")
}
