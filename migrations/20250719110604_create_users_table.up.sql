CREATE TABLE  IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    monthly_price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
