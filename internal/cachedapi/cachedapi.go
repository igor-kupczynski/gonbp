package cachedapi

import (
	"encoding/json"
	"gonbp/internal/nbpapi"
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

// Get returns the currency exchange rate for a given date from NBP table A
//
// Get first checks the on-disk cache and falls-back to nbpapi.Client.
func (c *Client) Get(curr string, day time.Time) (*nbpapi.Rates, error) {
	//TODO implement me
	panic("implement me")
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
