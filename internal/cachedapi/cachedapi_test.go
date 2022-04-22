package cachedapi

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/igor-kupczynski/gonbp/internal/nbpapi"
	"github.com/shopspring/decimal"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestCacheGetSet(t *testing.T) {
	key := cacheKey{
		curr: "EUR",
		day:  time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC),
	}
	v := &cacheValue{Rates: &nbpapi.Rates{
		Table:    "A",
		Currency: "euro",
		Code:     "EUR",
		Rates: []nbpapi.DailyRate{
			{
				No:            "074/A/NBP/2022",
				EffectiveDate: "2022-04-15",
				Mid:           decimal.NewFromFloat(4.2865),
			},
		},
	}}

	base, err := ioutil.TempDir("", "gonbp-TestCacheGetSet")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(base)

	c := &Client{
		dir: base,
		api: nil,
	}

	t.Run("read non-existing file", func(t *testing.T) {
		got, err := c.get(key)
		if got != nil {
			t.Errorf("Expected nil value, got %v", got)
			return
		}
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected path error, got %v", err)
			return
		}
	})

	t.Run("set a value", func(t *testing.T) {
		err := c.set(key, v)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
			return
		}
	})

	t.Run("read value previously set", func(t *testing.T) {
		got, err := c.get(key)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
			return
		}
		if diff := cmp.Diff(v, got); diff != "" {
			t.Errorf("get() mismatch (-want +got):\n%s", diff)
			return
		}
	})

	t.Run("read different value", func(t *testing.T) {
		got, err := c.get(cacheKey{
			curr: key.curr,
			day:  key.day.AddDate(0, 0, -1),
		})
		if got != nil {
			t.Errorf("Expected nil value, got %v", got)
			return
		}
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			t.Errorf("Expected path error, got %v", err)
			return
		}
	})

}
