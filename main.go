package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	CharityFilePath      string   `short:"c" description:"path to charity CSV"`
	GovernmentGrantsPath []string `short:"g" description:"path to a Form 990 Part VIII line 1e CSV from Open990"`
}

func main() {
	var opts Options
	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			panic(err)
		}
	}
}
