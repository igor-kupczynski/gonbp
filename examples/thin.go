package examples

import (
	"fmt"
	"github.com/igor-kupczynski/gonbp/gonbp"
	"log"
)

// Usage examples of the thin client
func ThinClientExamples(currencies []string) {
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
			case gonbp.NbpApiError:
				fmt.Printf("Can't find exchnage rate for today for %s\n", code)
			default:
				log.Fatal(err)
			}
		} else {
			fmt.Println(result)
		}
	}

	day := "2016-10-06"
	fmt.Printf("\n--- Exchange rates for %s\n", day)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Day("A", code, day)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	from := "2016-09-26"
	to := "2016-09-30"
	fmt.Printf("\n--- Exchange between [%s, %s]\n", from, to)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.DateRange("A", code, from, to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
