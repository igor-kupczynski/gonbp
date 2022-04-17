package gonbp

import (
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"testing"
)

func TestRate(t *testing.T) {
	nbp := New()

	t.Run("USD happy case", func(t *testing.T) {
		// {"table":"A","currency":"dolar amerykański","code":"USD","rates":[{"no":"043/A/NBP/2022","effectiveDate":"2022-03-03","mid":4.3257}]}
		want := &Rate{
			TableNo: "043/A/NBP/2022",
			Day:     day(2022, 3, 3),
			Mid:     decimal.NewFromFloat(4.3257),
		}
		got, err := nbp.Rate(USD, day(2022, 3, 3))
		if err != nil {
			t.Errorf("Rate() error = %v, want no error", err)
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
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
		wantErr := ErrNoExchangeRateForGivenDay{Day: day(2022, 4, 17)}
		_, gotErr := nbp.Rate(EUR, day(2022, 4, 17))
		if gotErr != wantErr {
			t.Errorf("Rate() error = %v, want %v", gotErr, wantErr)
		}
	})

	t.Run("Non-existing currency", func(t *testing.T) {
		// 404 NotFound
		wantErr := ErrNoRatesForCurrency{Curr: "DOGE"}
		_, gotErr := nbp.Rate("DOGE", day(2022, 4, 15))
		if gotErr != wantErr {
			t.Errorf("Rate() error = %v, want %v", gotErr, wantErr)
		}
	})

}
