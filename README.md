# Go NBP

This is a go client for the NBP currency exchange rates API.

[Narodowy Bank Polski](http://www.nbp.pl/) or NBP is the central bank of the
Republic of Poland. Among other tasks it publishes the official exchange rate
of Złoty (the Polish currency) against other currencies. [NBP provides an
API to access the exchange rates](http://api.nbp.pl/en.html).

**Note** This is an early work in progress.

### Thin client

This is a thin wrapper around the API provided by NBP. Its aim is to provide 
a 1:1 implementation of the functions available in the NBP Currency Exchange
API.
 
#### Examples

See the [examples](./examples) directory.

- Queries for particular currency
    
    * Current exchange rate
        ```go
        result, err := gonbp.DefaultNbpClient.Current("A", "EUR")
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(result)
        ```
    
    * Series of latest N exchange rates
        ```go
        result, err := gonbp.DefaultNbpClient.Last("A", "EUR", 5)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(result)
        ```
        
    * Exchange rate of currency for today
        ```go
        result, err := gonbp.DefaultNbpClient.Today("A", "EUR")
        if err != nil {
            switch err.(type) {
            case gonbp.NbpApiError:
                fmt.Printf("Can't find exchnage rate for today for EUR\n")
            default:
                log.Fatal(err)
            }
        } else {
            fmt.Println(result)
        }
        ```
    
    * Exchange rate of currency for given day
        ```go
        result, err := gonbp.DefaultNbpClient.Day("A", "EUR", "2016-10-10")
        if err != nil {
              log.Fatal(err)
        }
        fmt.Println(result)
        ```
    
    * Exchange rate of currency between two dates
        ```go
        result, err := gonbp.DefaultNbpClient.DateRange("A", "USD", "2016-09-01", "2016-09-15")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
        ```


This project is intended to run as a library, to check the examples run it as
a program.

```
$ go run main.go
--- Current exchange rates
&{A euro EUR [{195/A/NBP/2016 2016-10-07 4.2853}]}
&{A dolar amerykański USD [{195/A/NBP/2016 2016-10-07 3.8505}]}
&{A funt szterling GBP [{195/A/NBP/2016 2016-10-07 4.7872}]}
&{A jen (Japonia) JPY [{195/A/NBP/2016 2016-10-07 0.03708}]}
&{A rupia indonezyjska IDR [{195/A/NBP/2016 2016-10-07 0.00029653}]}

--- Last 3 exchange rates
&{A euro EUR [{193/A/NBP/2016 2016-10-05 4.3014} {194/A/NBP/2016 2016-10-06 4.2974} {195/A/NBP/2016 2016-10-07 4.2853}]}
&{A dolar amerykański USD [{193/A/NBP/2016 2016-10-05 3.8307} {194/A/NBP/2016 2016-10-06 3.8405} {195/A/NBP/2016 2016-10-07 3.8505}]}
&{A funt szterling GBP [{193/A/NBP/2016 2016-10-05 4.8783} {194/A/NBP/2016 2016-10-06 4.8873} {195/A/NBP/2016 2016-10-07 4.7872}]}
&{A jen (Japonia) JPY [{193/A/NBP/2016 2016-10-05 0.03724} {194/A/NBP/2016 2016-10-06 0.03705} {195/A/NBP/2016 2016-10-07 0.03708}]}
&{A rupia indonezyjska IDR [{193/A/NBP/2016 2016-10-05 0.00029477} {194/A/NBP/2016 2016-10-06 0.00029568} {195/A/NBP/2016 2016-10-07 0.00029653}]}

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

* **Thin client**
    - [ ] Queries for complete tables
    - [X] Queries for particular currency
        * [X] Current exchange rate
        * [X] Series of latest N exchange rates
        * [X] Exchange rate of currency for today
        * [X] Exchange rate of currency for given day
        * [X] Exchange rate of currency between two dates
    - [ ] Queries for gold prices

## Disclaimer

I'm not associated in any way with NBP, neither is this package.

This code is licenced under Apache Licence 2.0 see [LICENSE](./LICENSE)
for details.
