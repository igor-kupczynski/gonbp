package main

import (
	"flag"
	"fmt"
	"gonbp"
	"log"
	"strings"
	"time"
)

func main() {
	previous := flag.Bool("p", false, "fetch rate for the previous work day")
	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		log.Fatalf("Currency is required, e.g. nbp eur")
	}
	curr := gonbp.Currency(strings.ToUpper(args[0]))

	var date time.Time = time.Now()
	var err error
	if len(args) > 1 {
		date, err = time.Parse("2006-01-02", args[1])
		if err != nil {
			log.Fatalf("Can't parse date: %v", err)
		}
	}

	nbp, err := gonbp.Default()
	if err != nil {
		log.Fatalf("Can't create nbp client: %v", err)
	}

	var rate *gonbp.Rate
	if *previous {
		rate, err = nbp.PreviousRate(curr, date)
	} else {
		rate, err = nbp.Rate(curr, date)
	}
	if err != nil {
		log.Fatalf("Can't fetch rates: %v", err)
	}

	fmt.Printf("Table No: %s\n", rate.TableNo)
	fmt.Printf("     Day: %s\n", rate.Day.Format("2006-01-02"))
	fmt.Printf("    Rate: %s\n", rate.Mid)

}
