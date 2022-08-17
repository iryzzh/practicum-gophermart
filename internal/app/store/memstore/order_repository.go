package memstore

import (
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"sort"
	"time"
)

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Update(order *model.Order) error {
	o.store.Lock()
	defer o.store.Unlock()

	for _, v := range o.store.orders {
		if v.ID == order.ID {
			*v = *order
		}
	}

	return nil
}

func (o *OrderRepository) Exists(number string, userID int) error {
	if exists, _ := o.FindByNumber(number); exists != nil {
		if exists.UserID != userID {
			return store.ErrOrderConflict
		}
		return store.ErrOrderAlreadyExists
	}

	return nil
}

func (o *OrderRepository) Create(order *model.Order) error {
	if err := o.Exists(order.Number, order.UserID); err != nil {
		return err
	}

	o.store.Lock()
	defer o.store.Unlock()

	order.ID = o.store.orderNextID + 1

	// testing:
	if order.UploadedAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		order.UploadedAt = model.Time{Time: time.Now()}
	}

	o.store.orders[o.store.orderNextID] = order
	o.store.orderNextID++

	return nil
}

func (o *OrderRepository) FindByNumber(number string) (*model.Order, error) {
	o.store.Lock()
	defer o.store.Unlock()

	for _, v := range o.store.orders {
		if v.Number == number && !v.Deleted {
			return v, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (o *OrderRepository) FindByID(id int) (*model.Order, error) {
	o.store.Lock()
	defer o.store.Unlock()

	for _, v := range o.store.orders {
		if v.ID == id && !v.Deleted {
			return v, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (o *OrderRepository) FindByUserID(userID int) (*model.Order, error) {
	o.store.Lock()
	defer o.store.Unlock()

	for _, v := range o.store.orders {
		if v.UserID == userID && !v.Deleted {
			return v, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (o *OrderRepository) GetByUserID(userID int) ([]*model.Order, error) {
	var orders []*model.Order

	o.store.Lock()

	for _, v := range o.store.orders {
		if v.UserID == userID && !v.Deleted {
			orders = append(orders, v)
		}
	}

	o.store.Unlock()

	if len(orders) == 0 {
		return nil, store.ErrRecordNotFound
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.String() > orders[j].UploadedAt.String()
	})

	return orders, nil
}

func (o *OrderRepository) Incomplete() ([]*model.Order, error) {
	var orders []*model.Order

	o.store.Lock()

	for _, v := range o.store.orders {
		if v.Status == model.OrderNew || v.Status == model.OrderProcessing {
			orders = append(orders, v)
		}
	}

	o.store.Unlock()

	if len(orders) == 0 {
		return nil, store.ErrRecordNotFound
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.String() < orders[j].UploadedAt.String()
	})

	return orders, nil
}
