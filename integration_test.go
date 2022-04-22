package gonbp

import (
	"github.com/google/go-cmp/cmp"
	"github.com/igor-kupczynski/gonbp/internal/nbpapi"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestIntegrationRate(t *testing.T) {

	base, err := ioutil.TempDir("", "gonbp-integration test")
	if err != nil {
		t.Fatalf("Can't create the temp dir: %v", err)
		return
	}
	defer os.Remove(base)
	nbp := Init(base, http.DefaultClient)

	var nonCached, cached time.Duration

	t.Run("USD happy case", func(t *testing.T) {
		// {"table":"A","currency":"dolar amerykański","code":"USD","rates":[{"no":"043/A/NBP/2022","effectiveDate":"2022-03-03","mid":4.3257}]}
		want := &Rate{
			TableNo: "043/A/NBP/2022",
			Day:     day(2022, 3, 3),
			Mid:     decimal.NewFromFloat(4.3257),
		}
		t1 := time.Now()
		got, err := nbp.Rate(USD, day(2022, 3, 3))
		nonCached = time.Now().Sub(t1)
		if err != nil {
			t.Errorf("Rate() error = %v, want no error", err)
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("USD from cache", func(t *testing.T) {
		t1 := time.Now()
		_, err := nbp.Rate(USD, day(2022, 3, 3))
		cached = time.Now().Sub(t1)
		if err != nil {
			t.Errorf("Rate() error = %v, want no error", err)
			return
		}

		if cached > nonCached {
			t.Errorf(
				"Expected a faster response from cache, got cached %d ms vs non-cached %d ms",
				cached.Microseconds(),
				nonCached.Microseconds(),
			)
			return
		}
	})

	t.Run("EUR happy case", func(t *testing.T) {
		// {"table":"A","currency":"euro","code":"EUR","rates":[{"no":"021/A/NBP/2021","effectiveDate":"2021-02-02","mid":4.5025}]}
		want := &Rate{
			TableNo: "021/A/NBP/2021",
			Day:     day(2021, 2, 2),
			Mid:     decimal.NewFromFloat(4.5025),
		}
		got, err := nbp.Rate(EUR, day(2021, 2, 2))
		if err != nil {
			t.Errorf("Rate() error = %v, want no error", err)
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Not a working day", func(t *testing.T) {
		// 404 NotFound - Not Found - Brak danych
		wantErr := nbpapi.ErrNoExchangeRateForGivenDay
		_, gotErr := nbp.Rate(EUR, day(2022, 4, 17))
		if gotErr != wantErr {
			t.Errorf("Rate() error = %v, want %v", gotErr, wantErr)
		}
	})

	t.Run("Non-existing currency", func(t *testing.T) {
		// 404 NotFound
		wantErr := nbpapi.ErrNoRatesForCurrency
		_, gotErr := nbp.Rate("DOGE", day(2022, 4, 15))
		if gotErr != wantErr {
			t.Errorf("Rate() error = %v, want %v", gotErr, wantErr)
		}
	})
}

func TestIntegrationPreviousRate(t *testing.T) {
	base, err := ioutil.TempDir("", "gonbp-integration test")
	if err != nil {
		t.Fatalf("Can't create the temp dir: %v", err)
		return
	}
	defer os.Remove(base)
	nbp := Init(base, http.DefaultClient)

	t.Run("USD go back over a long weekend happy case", func(t *testing.T) {
		// {"table":"A","currency":"dolar amerykański","code":"USD","rates":[{"no":"074/A/NBP/2022","effectiveDate":"2022-04-15","mid":4.2865}]}
		want := &Rate{
			TableNo: "074/A/NBP/2022",
			Day:     day(2022, 4, 15),
			Mid:     decimal.NewFromFloat(4.2865),
		}
		got, err := nbp.PreviousRate(USD, day(2022, 4, 18))
		if err != nil {
			t.Errorf("Rate() error = %v, want no error", err)
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Non-existing currency", func(t *testing.T) {
		// 404 NotFound
		wantErr := nbpapi.ErrNoRatesForCurrency
		_, gotErr := nbp.Rate("DOGE", day(2022, 4, 16))
		if gotErr != wantErr {
			t.Errorf("Rate() error = %v, want %v", gotErr, wantErr)
		}
	})

}
