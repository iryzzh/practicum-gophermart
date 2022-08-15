package pgstore

import (
	"database/sql"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"sort"
	"time"
)

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Update(order *model.Order) error {
	stmt := `update orders set status = $1, accrual = $2, deleted = $3, deleted_at = $4 where number = $5;`

	res, err := o.store.db.Exec(stmt, order.Status.String(), order.Accrual, order.Deleted, order.DeletedAt, order.Number)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return store.ErrOrderUpdate
	}

	return nil
}

func (o *OrderRepository) Incomplete() ([]*model.Order, error) {
	var orders []*model.Order

	orderNew := model.OrderNew
	processing := model.OrderProcessing

	rows, err := o.store.db.Query(
		"select id, user_id, number, status, accrual, uploaded_at, deleted, deleted_at from orders where status in ($1,$2)",
		orderNew.String(), processing.String())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.Deleted, &order.DeletedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if len(orders) == 0 {
		return nil, store.ErrRecordNotFound
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.String() > orders[j].UploadedAt.String()
	})

	return orders, nil
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

	// testing:
	if order.UploadedAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		order.UploadedAt = model.Time{Time: time.Now()}
	}

	return o.store.db.QueryRow(
		"insert into orders (user_id, number, status, uploaded_at) values ($1, $2, $3, $4) returning id",
		order.UserID, order.Number, order.Status.String(), order.UploadedAt.String(),
	).Scan(&order.ID)
}

func (o *OrderRepository) FindByID(id int) (*model.Order, error) {
	order := &model.Order{}

	if err := o.store.db.QueryRow(
		"select id, user_id, number, status, accrual, uploaded_at from orders where id = $1",
		id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return order, nil
}

func (o *OrderRepository) FindByUserID(userID int) (*model.Order, error) {
	order := &model.Order{}

	if err := o.store.db.QueryRow(
		"select id, user_id, number, status, accrual, uploaded_at, deleted, deleted_at from orders where user_id = $1",
		userID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
		&order.Deleted,
		&order.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return order, nil
}

func (o *OrderRepository) FindByNumber(number string) (*model.Order, error) {
	order := &model.Order{}

	if err := o.store.db.QueryRow(
		"select id, number, user_id, status, accrual, uploaded_at, deleted, deleted_at from orders where number = $1",
		number,
	).Scan(
		&order.ID,
		&order.Number,
		&order.UserID,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
		&order.Deleted,
		&order.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return order, nil
}

func (o *OrderRepository) GetByUserID(userID int) ([]*model.Order, error) {
	var orders []*model.Order

	rows, err := o.store.db.Query(
		"select id, user_id, number, status, accrual, uploaded_at, deleted, deleted_at from orders where user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.Deleted, &order.DeletedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if len(orders) == 0 {
		return nil, store.ErrRecordNotFound
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.String() > orders[j].UploadedAt.String()
	})

	return orders, nil
}
