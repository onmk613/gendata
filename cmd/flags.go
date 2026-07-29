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
	flags.IntVar(&action.WriteConf.Concurrency, "concurrency", 1, "Number of concurrent workers to write data")
	flags.IntVar(&action.WriteConf.BatchSize, "batchsize", 1000, "Number of records to write in each batch")
	flags.IntVar(&action.WriteConf.RepeatCount, "repeatcount", 10, "Number of times to repeat the batch write (total records = batchsize * repeatcount)")
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
