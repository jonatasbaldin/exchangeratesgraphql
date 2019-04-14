package main

import (
	"bytes"
	"io/ioutil"

	"net/http"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

type fakeHttpClient struct{}

func (client fakeHttpClient) Get(url string) (*http.Response, error) {
	doc := etree.NewDocument()

	err := doc.ReadFromFile("fixtures/exchange-rate-sample.xml")
	if err != nil {
		return nil, err
	}

	body, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(body))}

	return resp, nil
}

func TestScrapper(t *testing.T) {
	var databaseURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := InitializeDB(databaseURL)
	if err != nil {
		t.Error(err)
	}

	e := &Env{db: db}
	e.migrateDB()
	e.clearDB()

	s := &Scrapper{db: db, httpClient: fakeHttpClient{}}

	err = s.Scrape()
	if err != nil {
		t.Error(err)
	}

	var exchangeRates []ExchangeRate
	db.Preload("Rates").Find(&exchangeRates)

	assert.Equal(t, len(exchangeRates), 5)

	// 32 from file + 1 EUR
	assert.Equal(t, len(exchangeRates[0].Rates), 33)

	e.clearDB()
}
