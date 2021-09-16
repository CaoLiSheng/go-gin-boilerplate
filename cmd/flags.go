package cmd

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pingcap/errors"
)

func (ahf allowHeaderFlag) String() string {
	return fmt.Sprintf("allow headers=%s", strings.Join(ahf, ";"))
}

func (ahf *allowHeaderFlag) Set(str string) error {
	*ahf = strings.Split(str, ";")
	return nil
}

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
	log.Println(flags.AllowHeaders)
	log.Println(flags.AllowOrigins)
	log.Println(flags.Dev)
	log.Println("DSN=", flags.DSN)
	log.Println("Super=", flags.SUN)
	log.Println("Password=", flags.SUP)
	log.Println(flags.Port)
}

var TheFlags = new(Flags)

func InitFlags() *Flags {
	TheFlags.AllowHeaders = allowHeaderFlag{"csrftoken", "session", "authorization"}
	flag.CommandLine.Var(&TheFlags.AllowHeaders, "headers", "allowed headers")

	TheFlags.AllowOrigins = allowOriginFlag{"http://localhost:4200", "http://localhost:3000", "http://localhost:3333"}
	flag.CommandLine.Var(&TheFlags.AllowOrigins, "origins", "allowed origins")

	TheFlags.Dev = devFlag(true)
	flag.CommandLine.Var(&TheFlags.Dev, "dev", "dev/release mode")

	flag.CommandLine.StringVar(&TheFlags.DSN, "dsn", "root:111111@tcp(127.0.0.1:3306)/test?charset=utf8mb4&collation=utf8mb4_bin&parseTime=true&loc=Local", "data source name")
	flag.CommandLine.StringVar(&TheFlags.SUN, "sun", "root", "super user name")
	flag.CommandLine.StringVar(&TheFlags.SUP, "sup", "111111", "super user password")

	TheFlags.Port = portFlag(9000)
	flag.CommandLine.Var(&TheFlags.Port, "port", "web server port")

	flag.Parse()
	
	TheFlags.Print()
	
	return TheFlags
}
