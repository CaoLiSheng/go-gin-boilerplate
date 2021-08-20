package cmd

type allowHeaderFlag []string

type allowOriginFlag []string

type devFlag bool

type portFlag int64

type Flags struct {
	AllowHeaders allowHeaderFlag
	AllowOrigins allowOriginFlag
	Dev          devFlag
	DSN          string
	SUN          string
	SUP          string
	Port         portFlag
}
