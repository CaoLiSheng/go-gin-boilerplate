package cmd

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pingcap/errors"
)

func (aof allowOriginFlag) String() string {
	return fmt.Sprintf("allow origins=%s", strings.Join(aof, ";"))
}

func (aof *allowOriginFlag) Set(str string) error {
	*aof = strings.Split(str, ";")
	return nil
}

func (df devFlag) String() string {
	if df {
		return "dev mode"
	} else {
		return "release mode"
	}
}

func (df *devFlag) Set(str string) error {
	switch str {
	case "true":
		*df = true
	case "false":
		*df = false
	default:
		return errors.New("only 'true'/'false' is allowed")
	}
	return nil
}

func (pf portFlag) String() string {
	return fmt.Sprintf("port=%d", pf)
}

func (pf *portFlag) Set(str string) error {
	i64, err := strconv.ParseInt(str, 10, 64)
	*pf = portFlag(i64)
	return err
}

func (flags *Flags) Print() {
	log.Println(flags.AllowOrigins)
	log.Println(flags.Dev)
	log.Println("DSN=", flags.DSN)
	log.Println("Super=", flags.SUN)
	log.Println("Password=", flags.SUP)
	log.Println(flags.Port)
}

func InitFlags() *Flags {
	flags := new(Flags)

	flags.AllowOrigins = allowOriginFlag{"http://localhost:4200", "http://localhost:3000", "http://localhost:3333"}
	flag.CommandLine.Var(&flags.AllowOrigins, "origins", "allowed origins")

	flags.Dev = devFlag(true)
	flag.CommandLine.Var(&flags.Dev, "dev", "dev/release mode")

	flag.CommandLine.StringVar(&flags.DSN, "dsn", "root:123456@tcp(192.168.1.6:3306)/test?charset=utf8&collation=utf8_bin&parseTime=true&loc=Local", "data source name")
	flag.CommandLine.StringVar(&flags.SUN, "sun", "root", "super user name")
	flag.CommandLine.StringVar(&flags.SUP, "sup", "111111", "super user password")

	flags.Port = portFlag(9000)
	flag.CommandLine.Var(&flags.Port, "port", "web server port")

	flag.Parse()
	
	flags.Print()
	
	return flags
}
