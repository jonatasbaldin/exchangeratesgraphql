package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shopspring/decimal"
)

const baseCurrency string = "EUR"

type ExchangeRate struct {
	gorm.Model
	Date  string `json:"date" gorm:"unique"`
	Rates []Rate `json:"rates"`
}

type Rate struct {
	gorm.Model
	Symbol         string          `json:"symbol"`
	Value          decimal.Decimal `json:"value" gorm:"type:numeric"`
	ExchangeRateID uint            `gorm:"index"`
}

func LatestRates(db *gorm.DB, base *string, symbols []*string) (latestRates []Rate, err error) {
	var exchangeRate ExchangeRate

	db.Preload("Rates").Last(&exchangeRate)

	err = CalculateRates(exchangeRate.Rates, base)
	if err != nil {
		return nil, err
	}

	latestRates = FilterSymbols(exchangeRate.Rates, symbols)

	return
}

func DatedRates(db *gorm.DB, base *string, symbols []*string, date string) (datedRates []Rate, err error) {
	var exchangeRate ExchangeRate

	db.Preload("Rates").Where(ExchangeRate{Date: date}).First(&exchangeRate)

	err = CalculateRates(exchangeRate.Rates, base)
	if err != nil {
		return nil, err
	}

	datedRates = FilterSymbols(exchangeRate.Rates, symbols)

	return
}

func HistoryExchangeRates(db *gorm.DB, base *string, symbols []*string, startAt string, endAt string) (historyExchangeRates []ExchangeRate, err error) {
	db.Preload("Rates").Where("date BETWEEN ? AND ?", startAt, endAt).Find(&historyExchangeRates)

	for i, er := range historyExchangeRates {
		err := CalculateRates(er.Rates, base)
		if err != nil {
			return nil, err
		}

		rates := FilterSymbols(er.Rates, symbols)

		historyExchangeRates[i].Rates = nil

		for _, c := range rates {
			historyExchangeRates[i].Rates = append(historyExchangeRates[i].Rates, c)
		}
	}

	return
}

func FilterSymbols(rates []Rate, symbols []*string) []Rate {
	if len(symbols) == 0 {
		return rates
	}

	var filteredRates []Rate

	for _, symbol := range symbols {
		for _, rate := range rates {
			if *symbol == rate.Symbol {
				filteredRates = append(filteredRates, rate)
			}
		}
	}

	return filteredRates
}

func CalculateRates(rates []Rate, base *string) error {
	if base != nil && *base != baseCurrency {
		baseRate, err := GetBaseRate(base, rates)
		if err != nil {
			return err
		}

		for i, c := range rates {
			rates[i].Value = c.Value.Div(baseRate)
		}
	}

	return nil
}

func GetBaseRate(base *string, rates []Rate) (baseRate decimal.Decimal, err error) {
	for _, rate := range rates {
		if rate.Symbol == *base {
			return rate.Value, nil
		}
	}

	return decimal.Zero, fmt.Errorf("base %s not supported", *base)
}

func InitializeDB(databaseURL string) (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", databaseURL)

	if err != nil {
		panic(err)
	}

	return
}
