package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"gendata/pkg/action"
	mydriver "gendata/pkg/driver"
)

func addRootFlags(flags *pflag.FlagSet) {
	flags.IntVar(&action.WriteConf.Concurrency, "concurrency", 1, "Number of concurrent workers to write data")
	flags.IntVar(&action.WriteConf.BatchSize, "batchsize", 1000, "Number of records to write in each batch")
	flags.IntVar(&action.WriteConf.RepeatCount, "repeatcount", 10, "Number of times to repeat the batch write (total records = batchsize * repeatcount)")

	var debug bool
	flags.BoolVar(&debug, "debug", false, "Print log")
	// 默认不输出日志
	if !debug {
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
	flags.StringVar(&mydriver.SqlConf.Table, "table", "", "Database table")
	flags.StringToStringVar(&mydriver.SqlConf.AdditionalArgs, "additionalargs", nil, "Additional connection arguments (key=value format, multiple args separated by comma)")
}
