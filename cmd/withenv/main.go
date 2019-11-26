package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/julian7/withenv"
)

var version = "SNAPSHOT"

func cmdName() string {
	return path.Base(os.Args[0])
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATH COMMAND [ARGS...]\n\n", cmdName())
	fmt.Fprint(
		os.Stderr,
		`Set environment variables defined in file PATH, and run COMMAND.

  -v       print program version
  -h       this help
Arguments:
  PATH     path to a file containing variable declarations in a KEY=VAL
           format. Spaces around KEY and VAL are NOT stripped. VAL has
		   variable expansion.
  COMMAND  next executable in line, which will be run with the newly set
           environment variables.
  ARGS...  any command line items are taken to COMMAND after variable
           expansion.
`)
}

func main() {
	flag.Usage = usage
	reqVersion := flag.Bool("v", false, "print program version")

	flag.Parse()

	if *reqVersion {
		fmt.Printf("%s version %v\n", cmdName(), version)
	}

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "PATH and COMMAND parameters are required.")
		usage()
		os.Exit(1)
	}

	env := withenv.New()
	if err := env.SetEnviron(os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "unable to import OS environment: %v\n", err)
		os.Exit(1)
	}

	err := env.Read(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading PATH: %v\n", err)
		os.Exit(1)
	}

	err = env.Run(args[1], args[2:]...)
	if err != nil {
		var xiterr *exec.ExitError
		if errors.As(err, &xiterr) {
			os.Exit(xiterr.ProcessState.ExitCode())
		}
	}
}
