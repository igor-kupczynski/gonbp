package main

import (
	"fmt"
	"log"
	"time"

	"github.com/igor-kupczynski/gonbp"
)

// ClientExamples shows how to use the gonbp client
func ClientExamples(currencies []string) {
	fmt.Println("\n--- Current exchange rates")
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Current("A", code)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	count := 3
	fmt.Printf("\n--- Last %d exchange rates\n", count)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Last("A", code, count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	fmt.Println("\n--- Exchange rates for today")
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Today("A", code)
		if err != nil {
			switch err.(type) {
			case gonbp.NbpAPIError:
				fmt.Printf("Can't find exchnage rate for today for %s\n", code)
			default:
				log.Fatal(err)
			}
		}
		fmt.Println(result)
	}

	dayStr := "2016-10-06"
	day, err := time.Parse(gonbp.DayFormat, dayStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n--- Exchange rates for %s\n", day)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Day("A", code, day)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	fromStr, toStr := "2016-09-26", "2016-09-30"
	from, err := time.Parse(gonbp.DayFormat, fromStr)
	if err != nil {
		log.Fatal(err)
	}
	to, err := time.Parse(gonbp.DayFormat, toStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n--- Exchange between [%s, %s]\n", from, to)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.DateRange("A", code, from, to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}

func main() {
	currencies := []string{"EUR", "USD", "GBP", "JPY", "IDR"}
	ClientExamples(currencies)
}
