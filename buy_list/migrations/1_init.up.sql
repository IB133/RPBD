BEGIN;

CREATE TABLE user(
    id SERIAL PRIMARY KEY,
    name varchar(255),
    chat_id int
);

CREATE TABLE fridge(
    user_id integer REFERENCES user(id) ON DELETE CASCADE,
    prod_name varchar(255),
    status boolean,
    experitation_date date,
);

CREATE TABLE buy_list(
    user_id integer REFERENCES user(id) ON DELETE CASCADE,
    prod_name varchar(255),
    weight real,
    buy_time tiemstamp
)

COMMIT;