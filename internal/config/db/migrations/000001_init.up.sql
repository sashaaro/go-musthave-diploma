CREATE SCHEMA IF NOT EXISTS gophermart;

CREATE TABLE gophermart.users (
  id SERIAL PRIMARY KEY,
  login TEXT NOT NULL users,
  password TEXT NOT NULL,
  wallet DECIMAL NOT NULL DEFAULT 0,
  withdrawn DECIMAL NOT NULL DEFAULT 0
);

CREATE TABLE gophermart.orders (
   id SERIAL PRIMARY KEY,
   number TEXT NOT NULL UNIQUE,
-- - `REGISTERED` — заказ зарегистрирован, но вознаграждение не рассчитано;
-- - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
-- - `PROCESSING` — расчёт начисления в процессе;
-- - `PROCESSED` — расчёт начисления окончен;
   status TEXT NOT NULL,
   accrual NUMERIC,
   uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
   user_id INTEGER REFERENCES gophermart.users(id)
);

CREATE TABLE gophermart.withdrawals (
    id SERIAL PRIMARY KEY,
    "order" TEXT NOT NULL,
    sum DECIMAL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES gophermart.users(id)
);
