CREATE TABLE person (
    id         bigserial PRIMARY KEY,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    email      text NOT NULL UNIQUE
);

---- create above / drop below ----

DROP TABLE person;