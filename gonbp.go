package gonbp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiRoot = "http://api.nbp.pl/api"

// DayFormat is the format in which NBP API expects dates
const DayFormat = "2006-01-02"

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
	Number        string
	EffectiveDate time.Time
	Mid           json.Number
}

// rawRates represent the rates as provided by the NBP api
type rawRates struct {
	Number        string      `json:"no"`
	EffectiveDate string      `json:"effectiveDate"`
	Mid           json.Number `json:"mid,Number"`
}

func (r *CurrencyRate) UnmarshalJSON(j []byte) error {
	var err error
	var raw rawRates
	if err = json.Unmarshal(j, &raw); err != nil {
		return err
	}

	if r.EffectiveDate, err = time.Parse(DayFormat, raw.EffectiveDate); err != nil {
		return err
	}
	r.Number = raw.Number
	r.Mid = raw.Mid
	return nil
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
	url := fmt.Sprintf("%s/%s/today", table, currCode)
	return c.fetchRates(url)
}

// Day is an exchange rate for a given day
func (c *NbpClient) Day(table string, currCode string, day time.Time) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s", table, currCode, timeToDay(day))
	return c.fetchRates(url)
}

// DateRange returns exchanges rates for the given date range
func (c *NbpClient) DateRange(table string, currCode string, fromDay time.Time, toDay time.Time) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", table, currCode, timeToDay(fromDay), timeToDay(toDay))
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

func timeToDay(t time.Time) string {
	return t.Format(DayFormat)
}
