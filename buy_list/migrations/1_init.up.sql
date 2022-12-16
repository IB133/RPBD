BEGIN;
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name varchar(255),
    chat_id int
);
CREATE TYPE status AS ENUM ('used', 'dispose', 'stored', 'opened');
CREATE TABLE fridge(
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    prod_name varchar(255),
    status status DEFAULT 'stored',
    experitation_date date,
    status_date date
);

CREATE TABLE buy_list(
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    prod_name varchar(255),
    weight real,
    buy_time timestamp
);
COMMIT;