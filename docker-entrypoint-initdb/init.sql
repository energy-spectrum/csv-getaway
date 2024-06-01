CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email varchar NOT NULL,
  name varchar NOT NULL,
  hashed_password varchar NOT NULL,
  role varchar NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ,
  avatar_url varchar
);

CREATE TABLE templates (
  id BIGSERIAL PRIMARY KEY,
  name varchar NOT NULL,
  filters json,
  skills json
);
