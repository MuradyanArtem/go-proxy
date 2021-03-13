package main

import (
	"flag"
	"os"
)

func main() {
	fs := flag.NewFlagSet("Proxy", flag.ExitOnError)

	var host string
	fs.StringVar(&host, "h", "", "print usage info")
	var port int
	fs.IntVar(&port, "p", 8000, "print build info")
	var config string
	fs.StringVar(&config, "c", "", "path to server config")

	SetupLogger(os.Stdout, "info")

	// nolint
	fs.Parse(os.Args[1:])

}
