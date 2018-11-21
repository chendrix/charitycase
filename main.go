package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"

	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	CharityFilePath              string `short:"c" description:"path to charity CSV" required:"true"`
	GovernmentGrantsDatabasePath string `short:"d" description:"path to a sqlite3 database populated with entries from Form 990 Part VIII line 1e CSV from Open990" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		} else {
			panic(err)
		}
	}

	errLog := log.New(os.Stderr, "", 0)
	database, err := sql.Open("sqlite3", opts.GovernmentGrantsDatabasePath)
	if err != nil {
		errLog.Printf("Failed to open sqlite database at path: %s\n%s\n", opts.GovernmentGrantsDatabasePath, err.Error())
		os.Exit(1)
	}

	row := database.QueryRow("SELECT COUNT(*) FROM grants")

	var count int
	err = row.Scan(&count)
	if err != nil {
		errLog.Printf("Failed to read entry count from database:\n%s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Grants Count: %d\n", count)

	c, err := os.Open(opts.CharityFilePath)
	if err != nil {
		errLog.Printf("Failed to open CSV file at path: %s\n%s\n", opts.CharityFilePath, err.Error())
		os.Exit(1)
	}
	defer c.Close()

	reader := csv.NewReader(bufio.NewReader(c))

	var charities []*CharityInput
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			errLog.Printf("Failed to read line: %s\n", err.Error())
			os.Exit(1)
		}
		ein := line[17]
		if ein == "" {
			continue
		}
		charities = append(charities, &CharityInput{
			BenevonInternalName: line[0],
			ExternalName:        line[1],
			Address: &Address{
				Line1:           line[2],
				Line2:           line[3],
				City:            line[4],
				StateOrProvince: line[5],
				Zip:             line[6],
			},
			Phone:       line[7],
			EIN:         ein,
			CreatedDate: line[15],
			CloseDate:   line[16],
		})
	}

	fmt.Printf("Number of charities with EINs: %d\n", len(charities))
}

type CharityInput struct {
	BenevonInternalName string
	ExternalName        string
	Address             *Address
	Phone               string
	EIN                 string
	CreatedDate         string
	CloseDate           string
}

type Address struct {
	Line1           string
	Line2           string
	City            string
	StateOrProvince string
	Zip             string
}

type CharityOutput struct {
	Charity CharityInput
	Grants  Grants
}

type Grants struct {
	_2004 string
	_2005 string
	_2006 string
	_2007 string
	_2008 string
	_2009 string
	_2010 string
	_2011 string
	_2012 string
	_2013 string
	_2014 string
	_2015 string
	_2016 string
	_2017 string
}

//
