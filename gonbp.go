package gonbp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiRoot = "http://api.nbp.pl/api"

// NbpAPIError encapsulates the errors returned by NBP API
type NbpAPIError struct {
	StatusCode int
	Message    string
}

func (c NbpAPIError) Error() string {
	return c.Message
}

// CurrencyRate represents an average currency rate published for a day
type CurrencyRate struct {
	Number        string      `json:"no"`
	EffectiveDate string      `json:"effectiveDate"`
	Mid           json.Number `json:"mid,Number"`
}

// CurrencyRateList of rates spanning for multiple days
type CurrencyRateList struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []CurrencyRate `json:"rates"`
}

// NbpClient is client to connect to the NBP exchange rate api
type NbpClient struct {
	// Http client to use to fetch the rate
	Client *http.Client

	// API root
	RatesAPIRoot string
}

// DefaultNbpClient is a client with default parameters
var DefaultNbpClient = &NbpClient{
	Client:       &http.Client{Timeout: time.Second * 5},
	RatesAPIRoot: apiRoot + "/exchangerates/rates",
}

// Current exchange rate for a currency
func (c *NbpClient) Current(table string, currCode string) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s", table, currCode)
	return c.fetchRates(url)
}

// Last n-exchange rates
func (c *NbpClient) Last(table string, currCode string, n int) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/last/%d", table, currCode, n)
	return c.fetchRates(url)
}

// Today exchange rate
func (c *NbpClient) Today(table string, currCode string) (*CurrencyRateList, error) {
	return c.Day(table, currCode, "today")
}

// Day is an exchange rate for a given day
func (c *NbpClient) Day(table string, currCode string, day string) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s", table, currCode, day)
	return c.fetchRates(url)
}

// DateRange returns exchanges rates for the given date range
func (c *NbpClient) DateRange(table string, currCode string, fromDay string, toDay string) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", table, currCode, fromDay, toDay)
	return c.fetchRates(url)
}

func (c *NbpClient) fetchRates(url string) (*CurrencyRateList, error) {
	response, err := c.Client.Get(c.RatesAPIRoot + "/" + url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, NbpAPIError{response.StatusCode, response.Status}
	}

	rates := &CurrencyRateList{}
	err = json.NewDecoder(response.Body).Decode(rates)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
