// Package gonbp is a go wrapper over NBP API
package gonbp

import (
	"errors"
	"fmt"
	"gonbp/internal/cachedapi"
	"gonbp/internal/nbpapi"
	"net/http"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/shopspring/decimal"
)

type nbpAPIClient interface {
	Get(curr string, day time.Time) (*nbpapi.Rates, error)
}

// NBP is the NBP API client
type NBP struct {
	api nbpAPIClient
}

// Init returns *NBP instance with a given httpClient
func Init(cacheDir string, client *http.Client) *NBP {
	return &NBP{api: cachedapi.Init(cacheDir, client)}
}

// Default returns *NBP instance using http.DefaultClient and $HOME/.config/nbp
func Default() (*NBP, error) {
	cacheDir, err := homedir.Expand("~/.config/nbp")
	if err != nil {
		return nil, err
	}
	return Init(cacheDir, http.DefaultClient), nil
}

// Currency enumerates supported currencies
type Currency string

const (
	CHF Currency = "CHF"
	EUR Currency = "EUR"
	USD Currency = "USD"
)

// Rate represents the currency exchange rate for a given date
type Rate struct {
	TableNo string
	Day     time.Time
	Mid     decimal.Decimal
}

// Rate returns the currency exchange rate for a given date from NBP table A
func (n *NBP) Rate(curr Currency, day time.Time) (*Rate, error) {
	apiRates, err := n.api.Get(string(curr), day)
	if err != nil {
		return nil, err
	}
	if len(apiRates.Rates) != 1 {
		return nil, fmt.Errorf("expectation failed: wanted a single rate, instead got %v", apiRates.Rates)
	}
	rate := apiRates.Rates[0]
	effectiveDay, err := time.Parse("2006-01-02", rate.EffectiveDate)
	if err != nil {
		return nil, fmt.Errorf("expectation failed: can't parse date as day %s", rate.EffectiveDate)
	}
	return &Rate{
		TableNo: rate.No,
		Day:     effectiveDay,
		Mid:     rate.Mid,
	}, nil
}

// PreviousRate returns the currency exchange rate for the last working day before the given day
func (n *NBP) PreviousRate(curr Currency, day time.Time) (*Rate, error) {
	checkForDay := day.AddDate(0, 0, -1)
	for {
		rate, err := n.Rate(curr, checkForDay)
		if errors.Is(err, nbpapi.ErrNoExchangeRateForGivenDay) {
			checkForDay = checkForDay.AddDate(0, 0, -1)
			continue
		}
		if err != nil {
			return nil, err
		}
		return rate, nil
	}
}
