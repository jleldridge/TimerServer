CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE hashed_passwords (
  id INT PRIMARY KEY,
  hashedSaltedPassword TEXT,
  FOREIGN KEY (id) REFERENCES users(id)
);

CREATE TABLE timer_entries (
  id SERIAL PRIMARY KEY,
  userId INT,
  startTime TIMESTAMP NOT NULL,
  stopTime TIMESTAMP,
  project TEXT,
  description TEXT,
  FOREIGN KEY (userId) REFERENCES users(id)
);