package main

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetBaseRateReturnsBase(t *testing.T) {
	value, _ := decimal.NewFromString("4.40")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
		},
	}

	base := "BRL"
	baseRate, err := GetBaseRate(&base, exchangeRate.Rates)

	assert.Equal(t, baseRate, value)
	assert.Nil(t, err)
}

func TestGetBaseRateBaseNotSupported(t *testing.T) {
	value, _ := decimal.NewFromString("4.40")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
		},
	}

	base := "ZAR"
	baseRate, err := GetBaseRate(&base, exchangeRate.Rates)

	assert.Equal(t, baseRate, decimal.Zero)
	assert.Equal(t, err.Error(), "base ZAR not supported")
}

func TestCalculateRatesWithSuccess(t *testing.T) {
	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	base := "USD"
	calculatedRate := exchangeRate.Rates[0].Value.Div(exchangeRate.Rates[1].Value)
	err := CalculateRates(exchangeRate.Rates, &base)

	assert.Nil(t, err)
	assert.Equal(t, exchangeRate.Rates[0].Value, calculatedRate)
}

func TestFilterSymbolsWithEmptySymbols(t *testing.T) {
	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	symbols := []*string{}
	filteredRates := FilterSymbols(exchangeRate.Rates, symbols)

	assert.Equal(t, len(filteredRates), 2)
}

func TestFilterSymbolsReturnsFilteredRates(t *testing.T) {
	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	usd := "USD"
	symbols := []*string{&usd}
	filteredRates := FilterSymbols(exchangeRate.Rates, symbols)

	assert.Equal(t, len(filteredRates), 1)
}

func TestLatestRatesReturnsLatestRate(t *testing.T) {
	e := InitializeTestEnv()
	e.migrateDB()
	e.clearDB()

	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	e.db.Create(&exchangeRate)

	// Create an older ExchangeRate to evalute if the latest one is returned
	exchangeRateOlder := ExchangeRate{
		Date: "2018-01-01",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	e.db.Create(&exchangeRateOlder)

	usd := "USD"
	brl := "BRL"
	symbols := []*string{&brl}
	calculatedValue := value.Div(valueUSD)

	latestRates, _ := LatestRates(e.db, &usd, symbols)

	assert.Equal(t, len(latestRates), 1)
	assert.Equal(t, latestRates[0].Value, calculatedValue)

	e.clearDB()
}
func TestDatedRatesReturnsDatedRate(t *testing.T) {
	e := InitializeTestEnv()
	e.migrateDB()
	e.clearDB()

	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-10",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	e.db.Create(&exchangeRate)

	valueOlder, _ := decimal.NewFromString("2.40")
	valueOlderUSD, _ := decimal.NewFromString("0.13")
	exchangeRateOlder := ExchangeRate{
		Date: "2018-01-09",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  valueOlder,
			},
			{
				Symbol: "USD",
				Value:  valueOlderUSD,
			},
		},
	}

	e.db.Create(&exchangeRateOlder)

	var base *string
	var symbols []*string

	datedRates, _ := DatedRates(e.db, base, symbols, "2018-01-09")

	assert.Equal(t, len(datedRates), 2)

	for _, r := range datedRates {

		switch r.Symbol {
		case "BRL":
			assert.Equal(t, r.Value, valueOlder)
		case "USD":
			assert.Equal(t, r.Value, valueOlderUSD)
		}
	}

	e.clearDB()
}

func TestDatedRatesInvalidDateReturnsEmpty(t *testing.T) {
	e := InitializeTestEnv()
	e.migrateDB()
	e.clearDB()

	var base *string
	var symbols []*string

	datedRates, _ := DatedRates(e.db, base, symbols, "2018-01-09")

	assert.Equal(t, len(datedRates), 0)

	e.clearDB()
}

func TestHistoricExchangeRatesReturnsExchangeRates(t *testing.T) {
	e := InitializeTestEnv()
	e.migrateDB()
	e.clearDB()

	value, _ := decimal.NewFromString("4.40")
	valueUSD, _ := decimal.NewFromString("1.13")
	exchangeRate := ExchangeRate{
		Date: "2019-01-10",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  value,
			},
			{
				Symbol: "USD",
				Value:  valueUSD,
			},
		},
	}

	e.db.Create(&exchangeRate)

	valueOlder, _ := decimal.NewFromString("2.40")
	valueOlderUSD, _ := decimal.NewFromString("0.13")
	exchangeRateOlder := ExchangeRate{
		Date: "2019-01-09",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  valueOlder,
			},
			{
				Symbol: "USD",
				Value:  valueOlderUSD,
			},
		},
	}

	e.db.Create(&exchangeRateOlder)

	// Create an ExchangeRate that shouldn't be returned
	exchangeRateMoreOlder := ExchangeRate{
		Date: "2018-01-09",
		Rates: []Rate{
			{
				Symbol: "BRL",
				Value:  valueOlder,
			},
			{
				Symbol: "USD",
				Value:  valueOlderUSD,
			},
		},
	}

	e.db.Create(&exchangeRateMoreOlder)

	var base *string
	var symbols []*string

	historyExchangeRates, _ := HistoryExchangeRates(e.db, base, symbols, "2019-01-09", "2019-01-10")

	assert.Equal(t, len(historyExchangeRates), 2)
}
