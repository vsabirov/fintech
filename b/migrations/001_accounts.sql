-- Write your migrate up statements here

CREATE TABLE accounts (
  id text PRIMARY KEY CHECK (id != '') NOT NULL,
  total double precision NOT NULL
);

---- create above / drop below ----

DROP TABLE accounts;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
