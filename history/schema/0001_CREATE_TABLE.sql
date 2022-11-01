CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  order_id INT,
  user_id INT,
  sum DECIMAL(18,2) CHECK ( sum > 0 ) NOT NULL,
  created TIMESTAMPTZ NOT NULL,
  modified TIMESTAMPTZ NOT NULL DEFAULT now(),
  status_id INT,
  FOREIGN KEY (status_id) REFERENCES status (status_id)
);

CREATE TABLE status (
    status_id SERIAL PRIMARY KEY,
    status_name VARCHAR(20)
);

CREATE TABLE order_has_position (
    order_id INT,
    position_id INT,
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (position_id) REFERENCES positions (id)
);

CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    good_id VARCHAR(255),
    amount INT CHECK ( amount > 0 ) NOT NULL,
    FOREIGN KEY (good_id) REFERENCES goods (good_id)
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    good_id VARCHAR(255),
    name VARCHAR(255),
    price DECIMAL(18,2) CHECK ( price > 0 ) NOT NULL,
    category VARCHAR(55),
    status_id INT,
    FOREIGN KEY (status_id) REFERENCES status (status_id)
);

INSERT INTO status (status_name) VALUES ('active'), ('archived'), ('done');