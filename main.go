/*
Copyright 2016 Igor Kupczynski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/igor-kupczynski/gonbp/gonbp"
	"fmt"
	"log"
)

func main() {
	currencies := []string{"EUR", "USD", "GBP", "JPY", "IDR"}


	fmt.Println("--- Current exchange rates")
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Current("A", code)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	count := 3
	fmt.Printf("--- Last %d exchange rates\n", count)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Last("A", code, count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	fmt.Println("--- Exchange rates for today\n")
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
	fmt.Printf("--- Exchange rates for %s\n", day)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.Day("A", code, day)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	from := "2016-09-26"
	to   := "2016-09-30"
 	fmt.Printf("--- Exchange between [%s, %s]\n", from, to)
	for _, code := range currencies {
		result, err := gonbp.DefaultNbpClient.DateRange("A", code, from, to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
