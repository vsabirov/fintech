-- Write your migrate up statements here

CREATE TABLE transfers (
  id text PRIMARY KEY CHECK (id != '') NOT NULL,
  amount double precision NOT NULL,

  sender text CHECK (sender != '') NOT NULL,
  receiver text CHECK (receiver != '') NOT NULL,

  FOREIGN KEY (sender) REFERENCES accounts (id),
  FOREIGN KEY (receiver) REFERENCES accounts (id)
)

---- create above / drop below ----

DROP TABLE transfers;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
