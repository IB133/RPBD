BEGIN;

CREATE TYPE role AS ENUM ('user', 'moder');

CREATE TABLE person(
    id SERIAL PRIMARY KEY,
    password varchar(255),
    login varchar(255),
    role role DEFAULT 'user'
);

CREATE TABLE users(
	id SERIAL PRIMARY KEY,
    email varchar(255)
)INHERITS (person);

CREATE TABLE news(
    id SERIAL PRIMARY KEY,
    theme varchar(255),
    title varchar(255),
    content text,
    posted_at date,
    moder_comm varchar(255),
    posted boolean DEFAULT false,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    moder_id integer REFERENCES person(id) ON DELETE CASCADE
);

CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    news_id integer REFERENCES news(id) ON DELETE CASCADE,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp,
    reply_to integer REFERENCES comments(id) ON DELETE CASCADE,
    content text,
    status bool DEFAULT true
);


COMMIT;