package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

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

	fmt.Fprintf(os.Stderr, "Grants Count: %d\n", count)

	c, err := os.Open(opts.CharityFilePath)
	if err != nil {
		errLog.Printf("Failed to open CSV file at path: %s\n%s\n", opts.CharityFilePath, err.Error())
		os.Exit(1)
	}
	defer c.Close()

	reader := csv.NewReader(bufio.NewReader(c))

	var validEIN = regexp.MustCompile(`\d+`)
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
		if !validEIN.MatchString(ein) {
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

	fmt.Fprintf(os.Stderr, "Number of charities with Valid EINs: %d\n", len(charities))

	var charitiesWithGrants []*CharityOutput

	for _, c := range charities {
		grants, err := SelectGrants(database, c.EIN)
		if err != nil {
			errLog.Println(err)

		}
		charitiesWithGrants = append(charitiesWithGrants, &CharityOutput{
			Charity: c,
			Grants:  grants,
		})
	}

	w := csv.NewWriter(os.Stdout)

	if err := w.Write(Header); err != nil {
		errLog.Println(err)
		os.Exit(1)
	}

	for _, c := range charitiesWithGrants {
		if err := w.Write(c.ToFormattedSlice()); err != nil {
			errLog.Println(err)
			os.Exit(1)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
}

var Header = []string{
	"EIN",
	"Benevon Name",
	"External Name",
	"Addr 1",
	"Addr 2",
	"City",
	"State/Province",
	"Zip",
	"Phone",
	"Created Date",
	"Close Date",
	"2004",
	"2005",
	"2006",
	"2007",
	"2008",
	"2009",
	"2010",
	"2011",
	"2012",
	"2013",
	"2014",
	"2015",
	"2016",
	"2017",
}

func (c *CharityOutput) ToFormattedSlice() []string {
	return []string{
		c.Charity.EIN,
		c.Charity.BenevonInternalName,
		c.Charity.ExternalName,
		c.Charity.Address.Line1,
		c.Charity.Address.Line2,
		c.Charity.Address.City,
		c.Charity.Address.StateOrProvince,
		c.Charity.Address.Zip,
		c.Charity.Phone,
		c.Charity.CreatedDate,
		c.Charity.CloseDate,
		c.Grants._2004,
		c.Grants._2005,
		c.Grants._2006,
		c.Grants._2007,
		c.Grants._2008,
		c.Grants._2009,
		c.Grants._2010,
		c.Grants._2011,
		c.Grants._2012,
		c.Grants._2013,
		c.Grants._2014,
		c.Grants._2015,
		c.Grants._2016,
		c.Grants._2017,
	}
}

func SelectGrants(db *sql.DB, ein string) (*Grants, error) {
	rows, err := db.Query("SELECT tax_period, value FROM grants WHERE ein = ?", ein)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("reading grants for EIN %s failed", ein))
	}

	var _2004, _2005, _2006, _2007, _2008, _2009, _2010, _2011, _2012, _2013, _2014, _2015, _2016, _2017 string
	for rows.Next() {
		var taxPeriod, value string
		rows.Scan(&taxPeriod, &value)

		year := string([]rune(taxPeriod)[0:4])

		switch year {
		case "2004":
			_2004 = value
		case "2005":
			_2005 = value
		case "2006":
			_2006 = value
		case "2007":
			_2007 = value
		case "2008":
			_2008 = value
		case "2009":
			_2009 = value
		case "2010":
			_2010 = value
		case "2011":
			_2011 = value
		case "2012":
			_2012 = value
		case "2013":
			_2013 = value
		case "2014":
			_2014 = value
		case "2015":
			_2015 = value
		case "2016":
			_2016 = value
		case "2017":
			_2017 = value
		default:
			return nil, fmt.Errorf("unable to match tax_period (%s) to year (%s) for ein (%s)", taxPeriod, year, ein)
		}
	}

	g := &Grants{
		_2004: _2004,
		_2005: _2005,
		_2006: _2006,
		_2007: _2007,
		_2008: _2008,
		_2009: _2009,
		_2010: _2010,
		_2011: _2011,
		_2012: _2012,
		_2013: _2013,
		_2014: _2014,
		_2015: _2015,
		_2016: _2016,
		_2017: _2017,
	}

	return g, nil
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
	Charity *CharityInput
	Grants  *Grants
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
