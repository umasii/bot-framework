package flags

import (
	"flag"
	"fmt"
	"os"
)

var Version = "v0.0.1"

const help = `
CicadaAIO`

var (
	FlagHelp    = flag.Bool("help", false, "print usage and help, and exit")
	FlagVersion = flag.Bool("version", false, "print version and exit")
	FlagDebug   = flag.Bool("v", false, "verbose mode")
	FlagRun     = flag.String("r", "", "run a test")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of framework:\n")
	fmt.Fprintf(os.Stderr, "\tframework [-v] [-r]\n")
	if *FlagHelp {
		fmt.Fprintln(os.Stderr, help)
	}
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}
