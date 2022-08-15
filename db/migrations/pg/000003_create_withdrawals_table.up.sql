create table if not exists withdrawals
(
    id           serial primary key,
    user_id      int       not null,
    order_number varchar   not null,
    withdraw     decimal(10, 2) default 0.0,
    processed_at timestamp not null
);

create unique index idx_withdrawals_user_id_order_number on withdrawals (user_id, order_number);
create unique index idx_withdrawals_user_id on withdrawals (user_id);