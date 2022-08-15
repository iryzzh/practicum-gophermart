package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"log"
	"net/http"
	"time"
)

var (
	ErrInternalError   = errors.New("internal error")
	ErrTooManyRequests = errors.New("too many requests")
)

type Accrual struct {
	address  string
	interval time.Duration
	store    store.Store
}

func New(store store.Store, address string, interval time.Duration) *Accrual {
	return &Accrual{
		store:    store,
		address:  address,
		interval: interval,
	}
}

func (a *Accrual) Serve(ctx context.Context) error {
	log.Println("Starting Accrual Client on", a.address)

	errChan := make(chan error, 1)
	go func() {
		select {
		case errChan <- a.processOrders(ctx):
		case <-ctx.Done():
		}
	}()

	var err error
	select {
	case err = <-errChan:
		log.Println("Error:", err)
	case <-ctx.Done():
		log.Println("Shutting down Accrual Client")
	}

	return err
}

func (a *Accrual) processOrders(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("runtime panic:", r)
		}
	}()

	ticker := time.NewTicker(a.interval)

	for {
		select {
		case <-ticker.C:
			orders, err := a.store.Order().Incomplete()
			if err != nil {
				if errors.Is(err, store.ErrRecordNotFound) {
					continue
				}
				return err
			}

			for _, order := range orders {
				if err := a.accrualOrderInfo(order); err != nil {
					log.Println("accrualOrderInfo err:", err)
					continue
				}

				if err := a.store.Order().Update(order); err != nil {
					log.Println("order update error:", err)
					return err
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (a *Accrual) accrualOrderInfo(order *model.Order) error {
	url := fmt.Sprintf("%s/api/orders/%s", a.address, order.Number)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return ErrInternalError
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}

	return nil
}
