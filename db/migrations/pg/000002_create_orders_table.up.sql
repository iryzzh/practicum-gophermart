create type STATUS as enum ('NEW','PROCESSING', 'INVALID', 'PROCESSED');

create table if not exists orders
(
    id          serial primary key,
    user_id     int            not null,
    number      varchar unique not null,
    status      STATUS         not null,
    accrual     decimal(10, 2) default 0.0,
    uploaded_at timestamp      not null,
    deleted     bool           default false,
    deleted_at  timestamp
);

create unique index idx_orders_number on orders (number);