package driver

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

var (
	SqlConf SqlConfiguration
)

type SqlConfiguration struct {
	Host           string
	Port           int
	User           string
	Password       string
	DBName         string
	AdditionalArgs map[string]string
}

// mysql dsn
func (c *SqlConfiguration) getMysqlDsn() string {
	// 设置默认值
	if c.User == "" {
		c.User = "root"
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.DBName == "" {
		c.DBName = "mysql"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	cfg := mysqlDriver.NewConfig()
	cfg.User = c.User
	cfg.Passwd = c.Password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	cfg.DBName = c.DBName
	cfg.Params = mergeParams(c.AdditionalArgs, "charset=utf8")

	return cfg.FormatDSN()
}

// pg dsn
func (c *SqlConfiguration) getPostgresDsn() string {
	// 设置默认值
	if c.User == "" {
		c.User = "postgres"
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.DBName == "" {
		c.DBName = "postgres"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   "/" + c.DBName,
	}

	q := u.Query()
	for key, value := range mergeParams(c.AdditionalArgs, "sslmode=disable") {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// ck dsn
func (c *SqlConfiguration) getClickhouseDsn() string {
	// 设置默认值
	if c.User == "" {
		c.User = "default"
	}
	if c.Port == 0 {
		c.Port = 9000
	}
	if c.DBName == "" {
		c.DBName = "default"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	u := &url.URL{
		Scheme: "tcp",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   "/" + c.DBName,
	}
	q := u.Query()
	for key, value := range mergeParams(c.AdditionalArgs, "write_timeout=20") {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// mergeParams 把默认参数与用户参数合并，用户参数优先。
func mergeParams(args map[string]string, defaultsArgs string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Fields(defaultsArgs)
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}
	for key, value := range args {
		params[key] = value
	}
	return params
}
