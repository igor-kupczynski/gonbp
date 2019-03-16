# Go NBP

This is a go client for the NBP currency exchange rates API.

[Narodowy Bank Polski](http://www.nbp.pl/) or NBP is the central bank of the
Republic of Poland. Among other tasks it publishes the official exchange rate
of Złoty (the Polish currency) against other currencies. [NBP provides an
API to access the exchange rates](http://api.nbp.pl/en.html).

**Note** This is a work in progress.

### Client

This is a thin wrapper around the API provided by NBP. Its aim is to provide 
a 1:1 implementation of the functions available in the NBP Currency Exchange
API.
 
#### Examples

See the [examples](./examples) directory.

- Queries for particular currency
    
    * Current exchange rate
        ```
		result, err := gonbp.DefaultNbpClient.Current("A", "EUR")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
        ```
    
    * Series of latest N exchange rates
        ```
		result, err := gonbp.DefaultNbpClient.Last("A", "EUR", 5)		
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
        ```
        
    * Exchange rate of currency for today
        ```
		result, err := gonbp.DefaultNbpClient.Today("A", "EUR")
		if err != nil {
			switch err.(type) {
			case gonbp.NbpAPIError:
				fmt.Printf("Can't find exchnage rate for today for %s\n", code)
			default:
				log.Fatal(err)
			}
		}
		fmt.Println(result)
        ```
    
    * Exchange rate of currency for given day
        ```
        day, err := time.Parse(gonbp.DayFormat, "2016-10-06")
        if err != nil {
            log.Fatal(err)
        }
        result, err := gonbp.DefaultNbpClient.Day("A", "EUR", day)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(result)
        ```
    
    * Exchange rate of currency between two dates
        ```
        from, err := time.Parse(gonbp.DayFormat, "2016-09-26")
        if err != nil {
            log.Fatal(err)
        }
        to, err := time.Parse(gonbp.DayFormat, "2016-09-30")
        if err != nil {
            log.Fatal(err)
        }
        result, err := gonbp.DefaultNbpClient.DateRange("A", code, from, to)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(result)
        ```


This project is intended to be included as a library, but you can run the
examples as a command.

```sh
$ cd examples/
examples$ go run .

--- Current exchange rates
&{A euro EUR [{053/A/NBP/2019 2019-03-15 4.3037}]}
&{A dolar amerykański USD [{053/A/NBP/2019 2019-03-15 3.8014}]}
&{A funt szterling GBP [{053/A/NBP/2019 2019-03-15 5.0406}]}
&{A jen (Japonia) JPY [{053/A/NBP/2019 2019-03-15 0.034044}]}
&{A rupia indonezyjska IDR [{053/A/NBP/2019 2019-03-15 0.00026657}]}

--- Last 3 exchange rates
&{A euro EUR [{051/A/NBP/2019 2019-03-13 4.3006} {052/A/NBP/2019 2019-03-14 4.3015} {053/A/NBP/2019 2019-03-15 4.3037}]}
&{A dolar amerykański USD [{051/A/NBP/2019 2019-03-13 3.8077} {052/A/NBP/2019 2019-03-14 3.8018} {053/A/NBP/2019 2019-03-15 3.8014}]}
&{A funt szterling GBP [{051/A/NBP/2019 2019-03-13 5.0013} {052/A/NBP/2019 2019-03-14 5.0398} {053/A/NBP/2019 2019-03-15 5.0406}]}
&{A jen (Japonia) JPY [{051/A/NBP/2019 2019-03-13 0.034201} {052/A/NBP/2019 2019-03-14 0.03403} {053/A/NBP/2019 2019-03-15 0.034044}]}
&{A rupia indonezyjska IDR [{051/A/NBP/2019 2019-03-13 0.00026693} {052/A/NBP/2019 2019-03-14 0.00026648} {053/A/NBP/2019 2019-03-15 0.00026657}]}

--- Exchange rates for today
Can't find exchnage rate for today for EUR
Can't find exchnage rate for today for USD
Can't find exchnage rate for today for GBP
Can't find exchnage rate for today for JPY
Can't find exchnage rate for today for IDR

--- Exchange rates for 2016-10-06
&{A euro EUR [{194/A/NBP/2016 2016-10-06 4.2974}]}
&{A dolar amerykański USD [{194/A/NBP/2016 2016-10-06 3.8405}]}
&{A funt szterling GBP [{194/A/NBP/2016 2016-10-06 4.8873}]}
&{A jen (Japonia) JPY [{194/A/NBP/2016 2016-10-06 0.03705}]}
&{A rupia indonezyjska IDR [{194/A/NBP/2016 2016-10-06 0.00029568}]}

--- Exchange between [2016-09-26, 2016-09-30]
&{A euro EUR [{186/A/NBP/2016 2016-09-26 4.3075} {187/A/NBP/2016 2016-09-27 4.2972} {188/A/NBP/2016 2016-09-28 4.2918} {189/A/NBP/2016 2016-09-29 4.3014} {190/A/NBP/2016 2016-09-30 4.3120}]}
&{A dolar amerykański USD [{186/A/NBP/2016 2016-09-26 3.8324} {187/A/NBP/2016 2016-09-27 3.8227} {188/A/NBP/2016 2016-09-28 3.8264} {189/A/NBP/2016 2016-09-29 3.8354} {190/A/NBP/2016 2016-09-30 3.8558}]}
&{A funt szterling GBP [{186/A/NBP/2016 2016-09-26 4.9560} {187/A/NBP/2016 2016-09-27 4.9573} {188/A/NBP/2016 2016-09-28 4.9717} {189/A/NBP/2016 2016-09-29 4.9837} {190/A/NBP/2016 2016-09-30 4.9962}]}
&{A jen (Japonia) JPY [{186/A/NBP/2016 2016-09-26 0.038081} {187/A/NBP/2016 2016-09-27 0.038018} {188/A/NBP/2016 2016-09-28 0.037981} {189/A/NBP/2016 2016-09-29 0.037798} {190/A/NBP/2016 2016-09-30 0.038171}]}
&{A rupia indonezyjska IDR [{186/A/NBP/2016 2016-09-26 0.0002941} {187/A/NBP/2016 2016-09-27 0.00029508} {188/A/NBP/2016 2016-09-28 0.00029557} {189/A/NBP/2016 2016-09-29 0.00029561} {190/A/NBP/2016 2016-09-30 0.00029533}]}

```


## Implementation Status

* **Client**
    - [ ] Queries for complete tables
    - [X] Queries for particular currency
        * [X] Current exchange rate
        * [X] Series of latest N exchange rates
        * [X] Exchange rate of currency for today
        * [X] Exchange rate of currency for given day
        * [X] Exchange rate of currency between two dates
    - [ ] Queries for gold prices

## Reach out

Let me know if my go sucks, or suggest an improvements to this lib.

## Disclaimer

I'm not associated in any way with NBP, neither is this package.

This code is licenced under Apache Licence 2.0 see [LICENSE](./LICENSE)
for details.
