package cmd

type allowOriginFlag []string

type devFlag bool

type portFlag int64

type Flags struct {
	AllowOrigins allowOriginFlag
	Dev          devFlag
	DSN          string
	SUN          string
	SUP          string
	Port         portFlag
}
