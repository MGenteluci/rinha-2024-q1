CREATE TABLE clients (
  id SERIAL PRIMARY KEY,
  client_limit INT,
  balance INT
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  clientID INT FOREIGN KEY(clientID) REFERENCES clients(id),
  transaction_value INT
  transaction_type VARCHAR(1)
  transaction_description VARCHAR(10)
  transaction_date DATE
);

DO $$
BEGIN
  INSERT INTO clients (id, client_limit, balance)
  VALUES
    (1, 100000, 0),
    (2, 80000, 0),
    (3, 1000000, 0),
    (4, 10000000, 0),
    (5, 500000, 0);
END; $$
