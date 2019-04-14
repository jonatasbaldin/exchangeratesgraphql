package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/beevik/etree"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

const last90DaysRatesUrl string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
const historicalUrl string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml"

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Scrapper struct {
	httpClient httpClient
	db         *gorm.DB
}

func (s *Scrapper) Scrape() (err error) {
	var url string

	eurValue, err := decimal.NewFromString("1.0")
	if err != nil {
		return
	}

	eur := Rate{
		Symbol: "EUR",
		Value:  eurValue,
	}

	switch s.hasHistoricalData() {
	case true:
		url = last90DaysRatesUrl
	case false:
		url = historicalUrl
	}

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc := etree.NewDocument()
	err = doc.ReadFromBytes(body)
	if err != nil {
		return
	}

	for _, cube := range doc.FindElements(".//Cube/*[@time]") {
		date := cube.SelectAttr("time")
		if date == nil {
			return errors.New("date not found while scrapping")
		}

		exchangeRate := ExchangeRate{
			Date: date.Value,
		}

		cubes := cube.SelectElements("Cube")
		if len(cubes) == 0 {
			return fmt.Errorf("rates not found under date %s", date.Value)
		}

		for _, rate := range cubes {
			value := rate.SelectAttr("rate")
			if value == nil {
				return fmt.Errorf("value not found under rate %v", rate)
			}

			rateValue, err := decimal.NewFromString(value.Value)
			if err != nil {
				return err
			}

			symbol := rate.SelectAttr("currency")
			if symbol == nil {
				return fmt.Errorf("symbol not found under rate %v", rate)
			}

			rate := Rate{
				Symbol: symbol.Value,
				Value:  rateValue,
			}

			exchangeRate.Rates = append(exchangeRate.Rates, rate)
		}

		exchangeRate.Rates = append(exchangeRate.Rates, eur)

		err := s.db.Where(ExchangeRate{Date: exchangeRate.Date}).FirstOrCreate(&exchangeRate).Error

		if err != nil {
			log.WithFields(log.Fields{
				"exchangeRateDate": exchangeRate.Date,
			}).Info(err)
		}

		log.WithFields(log.Fields{
			"exchangeRateDate": exchangeRate.Date,
		}).Info("exchangeRate created successfully")
	}

	return
}

func (s *Scrapper) hasHistoricalData() bool {
	return s.db.Where(ExchangeRate{Date: "1999-01-04"}).RecordNotFound()
}
