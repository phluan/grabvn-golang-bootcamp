CREATE TABLE users (
        id uuid DEFAULT uuid_generate_v4(),
        email VARCHAR NOT NULL,
        password VARCHAR NOT NULL,
        first_name VARCHAR,
        last_name VARCHAR,
        PRIMARY KEY (id)
);
