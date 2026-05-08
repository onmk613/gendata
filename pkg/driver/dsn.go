package driver

import (
	"strconv"
	"strings"
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
	// 多个空格分隔的参数
	defaultsArgs := "charset=utf8"
	arg := additionalArgsToDsn(c.AdditionalArgs, defaultsArgs, "&")

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

	// dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	if c.Password == "" {
		return c.User + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.DBName + "?" + arg
	}
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.DBName + "?" + arg
}

// pg dsn
func (c *SqlConfiguration) getPostgresDsn() string {
	// 多个空格分隔的参数
	defaultsArgs := "sslmode=disable"
	arg := additionalArgsToDsn(c.AdditionalArgs, defaultsArgs, " ")

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

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	if c.Password == "" {
		return "user=" + c.User + " host=" + c.Host + " port=" + strconv.Itoa(c.Port) + " dbname=" + c.DBName + " " + arg
	}
	return "user=" + c.User + " password=" + c.Password + " host=" + c.Host + " port=" + strconv.Itoa(c.Port) + " dbname=" + c.DBName + " " + arg
}

// ck dsn
func (c *SqlConfiguration) getClickhouseDsn() string {
	// 多个空格分隔的参数
	defaultsArgs := "write_timeout=20"
	arg := additionalArgsToDsn(c.AdditionalArgs, defaultsArgs, "&")

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

	// dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
	if c.Password == "" {
		return "tcp://" + c.Host + ":" + strconv.Itoa(c.Port) + "?username=" + c.User + "&database=" + c.DBName + "&" + arg
	}
	return "tcp://" + c.Host + ":" + strconv.Itoa(c.Port) + "?username=" + c.User + "&password=" + c.Password + "&database=" + c.DBName + "&" + arg
}

// args
func additionalArgsToDsn(args map[string]string, defaultsArgs, symbol string) string {
	var arg string

	if args == nil {
		args = make(map[string]string)
	}

	pairs := strings.Fields(defaultsArgs)
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			if _, ok := args[parts[0]]; ok {
				continue
			}
			args[parts[0]] = parts[1]
		}
	}

	// 多一层判断
	if len(args) > 0 {
		var dsn string
		for key, value := range args {
			dsn += key + "=" + value + symbol
		}

		arg += dsn[:len(dsn)-1]
	}

	return arg
}
