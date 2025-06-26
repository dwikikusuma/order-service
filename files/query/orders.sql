CREATE TABLE order_detail (
	id BIGSERIAL PRIMARY KEY,
	products TEXT NOT NULL,
	order_history TEXT NOT NULL
);

CREATE TABLE orders(
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL,
	amount NUMERIC NOT NULL,
	total_qty INTEGER NOT NULL,
	status INTEGER NOT NULL,
	payment_method VARCHAR(50),
	shipping_address TEXT,
	order_detail_id BIGINT REFERENCES order_detail(id),
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_request_log(
	id BIGSERIAL PRIMARY KEY,
	idempotency_token TEXT UNIQUE NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

