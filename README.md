# Go NBP

This is a go client for the NBP currency exchange rates API.

[Narodowy Bank Polski](http://www.nbp.pl/) or NBP is the central bank of the
Republic of Poland. Among other tasks it publishes the official exchange rate
of Złoty (the Polish currency) against other currencies. [NBP provides an
API to access the exchange rates](http://api.nbp.pl/en.html).

## TODO

This is an early work in progress.

### Next

- more currency rate functions
- description of usage in the read me

## Implementation Status

- [ ] Queries for complete tables
- [ ] Queries for particular currency
  * [X] Current exchange rate
  * [X] Series of latest N exchange rates
  * [ ] Exchange rate of currency for today
  * [X] Exchange rate of currency for given day
  * [ ] Exchange rate of currency between two dates
- [ ] Queries for gold prices

## Disclaimer

I'm not associated in any way with NBP, neither is this package.

This code is licenced under Apache Licence 2.0 see [LICENSE.md](./LICENSE.md)
for details.
