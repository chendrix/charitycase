package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"

	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	CharityFilePath              string `short:"c" description:"path to charity CSV"`
	GovernmentGrantsDatabasePath string `short:"d" description:"path to a sqlite3 database populated with entries from Form 990 Part VIII line 1e CSV from Open990"`
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

	errLog := log.New(os.Stderr, "", 0)
	database, err := sql.Open("sqlite3", opts.GovernmentGrantsDatabasePath)
	if err != nil {
		errLog.Printf("Failed to open sqlite database at path: %s\n%s", opts.GovernmentGrantsDatabasePath, err.Error())
		os.Exit(1)
	}

	row := database.QueryRow("SELECT COUNT(*) FROM grants")

	var count int
	err = row.Scan(&count)
	if err != nil {
		errLog.Printf("Failed to read entry count from database:\n%s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Count: %d\n", count)
}

//
