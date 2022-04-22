# gonbp

This is a go wrapper over the NBP currency exchange rates API.

[Narodowy Bank Polski](http://www.nbp.pl/) (NBP) is the central bank of the
Republic of Poland. Among other tasks it publishes the official exchange rate
of _Złoty_ (the Polish currency) against other currencies. [NBP provides an
API to access the exchange rates](http://api.nbp.pl/en.html).

## Use from CLI

Fetch current day CHF rate 
```shell
❯ nbp CHF
Table No: 078/A/NBP/2022
     Day: 2022-04-22
    Rate: 4.493
```

Fetch previous day USD rate
```shell
❯ nbp -p USD
Table No: 077/A/NBP/2022
     Day: 2022-04-21
    Rate: 4.2596
```

Fetch EUR for a given day
```shell
❯ nbp EUR 2022-04-15
Table No: 074/A/NBP/2022
     Day: 2022-04-15
    Rate: 4.6378
```



## Use as a library

See [`integration_test.go`](https://github.com/igor-kupczynski/gonbp/blob/main/gonbp_test.go).
