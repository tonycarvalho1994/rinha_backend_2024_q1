CREATE DATABASE rinha2024q1;

\c rinha2024q1;

CREATE TABLE customers (
    id CHAR(1) PRIMARY KEY,
    limite FLOAT
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    customer_id char(1) REFERENCES customers(id),
    valor FLOAT,
    tipo CHAR(1),
    descricao CHAR(10),
    realizada_em TIMESTAMP,
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE INDEX idx_customer_id_realizada_em ON transactions(customer_id, realizada_em);

INSERT INTO customers(id, limite) VALUES ('1', 100000);
INSERT INTO customers(id, limite) VALUES ('2', 80000);
INSERT INTO customers(id, limite) VALUES ('3', 1000000);
INSERT INTO customers(id, limite) VALUES ('4', 10000000);
INSERT INTO customers(id, limite) VALUES ('5', 500000);
