package main

// var (
// 	admincfg   = servercfg.Export(flag.CommandLine, prefix+".admin.")
// 	sniffercfg = servercfg.Export(flag.CommandLine, prefix+".sniffer.")
// 	dbconf     = dbcfg.Export(flag.CommandLine, prefix+".postgres.")
// )

// func Parse(fs *flag.FlagSet, argv []string) {
// 	fs.Func("c", "path to config.yml `file`", func(path string) error {
// 		data, err := ioutil.ReadFile(path)
// 		if err != nil {
// 			return err
// 		}
// 		return flagutil.ParseYAML(fs, data)
// 	})
// 	var isHelp bool
// 	fs.BoolVar(&isHelp, "h", false, "print usage info")

// 	var logLevel string
// 	fs.StringVar(&logLevel, prefix+".log", "info", "set logging level")

// 	// nolint
// 	fs.Parse(argv)

// 	if isHelp {
// 		fs.PrintDefaults()
// 		os.Exit(0)
// 	}

// 	SetupLogger(os.Stdout, logLevel)
// }
