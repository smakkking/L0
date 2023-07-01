CREATE TABLE Delivery (
    del_id BIGINT PRIMARY KEY,

    del_name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    del_address TEXT,
    region TEXT,
    email TEXT,
);

CREATE TABLE Payment (
    payment_id BIGINT PRIMARY KEY,

    payment_transaction TEXT NOT NULL,
    request_id TEXT,
    currency VARCHAR(8),
    payment_provider VARCHAR(8),
    amount INT,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
);

CREATE TABLE Item (
    chrt_id BIGINT PRIMARY KEY,
    
    track_number TEXT,
    price INT,
    rid TEXT,
    item_name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id BIGINT,
    brand TEXT,
    item_status INT
);

CREATE TABLE Order (
    order_uid UUID PRIMARY KEY

    track_number TEXT,
    order_entry ,
    delivery_id BIGINT,
    payment_id BIGINT,

    FOREIGN KEY (delivery_id) REFERENCES Delivery(del_id)
    FOREIGN KEY (payment_id) REFERENCES Payment(payment_id)
);

CREATE TABLE 