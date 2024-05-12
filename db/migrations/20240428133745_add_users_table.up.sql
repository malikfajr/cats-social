CREATE TABLE IF NOT EXISTS "users" (
   id BIGSERIAL PRIMARY KEY,
   email VARCHAR (50) UNIQUE NOT NULL,
   name VARCHAR (50) NOT NULL,
   password CHAR (60) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_email ON users(email);