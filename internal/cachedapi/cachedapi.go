package cachedapi

import (
	"encoding/json"
	"errors"
	"gonbp/internal/nbpapi"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type nbpAPIClient interface {
	Get(curr string, day time.Time) (*nbpapi.Rates, error)
}

// Client is a low-level client over the NBP Rates API
type Client struct {
	dir string
	api nbpAPIClient
}

// Init returns *Rates instance with net/http.DefaultClient
func Init(cacheDir string, client *http.Client) *Client {
	return &Client{
		dir: cacheDir,
		api: nbpapi.Init(client),
	}
}

type cacheKey struct {
	curr string
	day  time.Time
}

func (k *cacheKey) dir() string {
	return k.curr
}

func (k *cacheKey) fname() string {
	return k.day.Format("2006-01-02.json")
}

type cacheValue struct {
	Rates *nbpapi.Rates `json:"Rates,omitempty"`
}

var noValueForDay = &cacheValue{Rates: nil}

// Get returns the currency exchange rate for a given date from NBP table A
//
// Get first checks the on-disk cache and falls-back to nbpapi.Client.
func (c *Client) Get(curr string, day time.Time) (*nbpapi.Rates, error) {
	key := cacheKey{curr: curr, day: day}
	v, err := c.get(key)

	if err == nil {
		if v.Rates == nil {
			return nil, nbpapi.ErrNoExchangeRateForGivenDay
		}
		return v.Rates, nil
	}

	var pathError *fs.PathError
	if !errors.As(err, &pathError) {
		return nil, err
	}

	got, err := c.api.Get(curr, day)
	if err == nbpapi.ErrNoExchangeRateForGivenDay {
		if err := c.set(key, noValueForDay); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err := c.set(key, &cacheValue{got}); err != nil {
		return nil, err
	}

	return got, nil
}

func (c *Client) get(k cacheKey) (*cacheValue, error) {
	dir := path.Join(c.dir, k.dir())
	buf, err := ioutil.ReadFile(path.Join(dir, k.fname()))
	if err != nil {
		return nil, err
	}
	var v cacheValue
	if err := json.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) set(k cacheKey, v *cacheValue) error {
	dir := path.Join(c.dir, k.dir())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dir, k.fname()), buf, 0644)
}
