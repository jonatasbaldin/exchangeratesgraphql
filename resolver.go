package main

import (
	"context"
)

type Resolver struct {
	Env Env
}

func (r *Resolver) Rate() RateResolver {
	return &rateResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type rateResolver struct{ *Resolver }

func (r *rateResolver) Value(ctx context.Context, obj *Rate) (string, error) {
	return obj.Value.String(), nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Latest(ctx context.Context, base *string, symbols []*string) ([]Rate, error) {
	latestRates, err := LatestRates(r.Env.db, base, symbols)

	if err != nil {
		return nil, err
	}

	return latestRates, nil
}

func (r *queryResolver) Date(ctx context.Context, base *string, symbols []*string, date string) ([]Rate, error) {
	datedRates, err := DatedRates(r.Env.db, base, symbols, date)

	if err != nil {
		return nil, err
	}

	return datedRates, nil
}

func (r *queryResolver) History(ctx context.Context, base *string, symbols []*string, startAt string, endAt string) ([]ExchangeRate, error) {
	historyExchangeRates, err := HistoryExchangeRates(r.Env.db, base, symbols, startAt, endAt)

	if err != nil {
		return nil, err
	}

	return historyExchangeRates, nil
}
