# Exchange Rates GraphQL

<p align="center">
  <a href="https://exchangeratesgraphql.deployeveryday.com"><img src="static/logo.png"></a>
</p>

Exchange Rates GraphQL is a free service for current and historical foreign exchange rates [published by the European Central Bank](https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html). _Inspired by [exchangeratesapi](https://github.com/exchangeratesapi/exchangeratesapi), inlcuding this description_.

## Queries

To get the latest Exchange Rate:
```graphql
query {
  latest {
    symbol
    value
  }
}
```

You can also ask for a specific base currency and which symbols to return. **All Queries accept these parameters**.
```graphql
query {
  latest(base: "USD", symbols: ["EUR", "USD", "BRL"]) {
    symbol
    value
  }
}
```

Getting a specific date:
```graphql
query {
  date(date: "2000-01-03") {
    symbol
    value
  }
}
```

Finally, asking for a date range:

```graphql
query {
  history(startAt: "2019-01-01", endAt: "2019-01-10") {
    rates {
      symbol
      value
    }
  }
}

```

## Developing
To run it locally you will need:
- Go v1.11 or higher
- PostgreSQL

First, set the following environment variables:
```bash
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'
export TEST_DATABASE_URL='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' # used only for tests!
export GIN_MODE='release'
export PORT=8000
```

Install the dependencies, build and run:
```bash
$ make build
$ make run
```

Go to your browser at `http://localhost:8080/graphql/playground` and have fun!

The currency data comes from the European Central Bank. The builtin scrapper will download the data from since _1999_ in its first run, and subsequently download the latest 30 days and update accordingly. To run it, use:
```bash
$ make scrape
```

To run the tests, use the following command:
```
$ make test
```

## FAQ
**Q: Why am a specific date doesn't return any data?**   
R: The source data on the European Central bank if for every working day, except TARGET closing days. Also, the rates are updated every hour.

## License
[MIT](./LICENSE).