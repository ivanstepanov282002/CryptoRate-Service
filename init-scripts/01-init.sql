-- This script is executed the FIRST time the container is launched
-- Create a table to store exchange rates
CREATE TABLE IF NOT EXISTS Users ( 
user_id BIGINT NOT NULL PRIMARY KEY, -- ID in Telegram 
user_name VARCHAR(100), 
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Settings ( 
id SERIAL PRIMARY KEY, 
user_id BIGINT UNIQUE NOT NULL, 
time_interval INTEGER, 
last_sent TIMESTAMP, 
FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS Currency ( 
id SERIAL PRIMARY KEY, 
name_currency VARCHAR(50) UNIQUE
);

CREATE TABLE IF NOT EXISTS Currency_settings ( 
user_id BIGINT NOT NULL, 
currency_id INTEGER NOT NULL, 
is_active BOOLEAN DEFAULT true, 
PRIMARY KEY (user_id, currency_id), 
FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE, 
FOREIGN KEY (currency_id) REFERENCES Currency(id) 
);

CREATE TABLE IF NOT EXISTS Exchange_rate ( 
id SERIAL PRIMARY KEY, 
currency_id INTEGER NOT NULL,
price DECIMAL(15, 6) NOT NULL, 
recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY  (currency_id) REFERENCES Currency(id)
);