CREATE TABLE "user" 
(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(20) NOT NULL,
    surname VARCHAR(20) NOT NULL,
    join_date DATE NOT NULL,
    hierarchy SMALLINT NOT NULL,
    activated BOOLEAN NOT NULL
);

CREATE TABLE "review" 
(
    id SERIAL PRIMARY KEY,
    by_user_fk INTEGER NOT NULL REFERENCES "user" (id),
    comment TEXT NOT NULL
);

CREATE TABLE "basic_auth" 
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    passwd VARCHAR(15) NOT NULL,
    user_fk INTEGER NOT NULL REFERENCES "user" (id)
);

CREATE TABLE "content" 
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    subtitle VARCHAR(500),
    body TEXT NOT NULL,
    author_user_fk INTEGER NOT NULL REFERENCES "user" (id)
);

CREATE TABLE "report" 
(
    id SERIAL PRIMARY KEY,
    by_user_fk INTEGER NOT NULL REFERENCES "user" (id),
    to_user_fk INTEGER NOT NULL REFERENCES "user" (id),
    comment TEXT NOT NULL,
    review_fk INTEGER REFERENCES "review" (id)
);

CREATE TABLE "article" 
(
    id SERIAL PRIMARY KEY,
    author_user_fk INTEGER NOT NULL REFERENCES "user" (id),
    publish_date DATE NOT NULL,
    content_fk INTEGER NOT NULL REFERENCES "content" (id),
    rating FLOAT NOT NULL DEFAULT 0.0,
    published BOOLEAN NOT NULL DEFAULT TRUE
);