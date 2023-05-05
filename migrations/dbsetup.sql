CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE problem (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    level VARCHAR(20) NOT NULL,
    samples VARCHAR(255),
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    CONSTRAINT problem_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE topic (
    id SERIAL PRIMARY KEY,
    topics VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE relation (
    id SERIAL PRIMARY KEY,
    problem_id INT NOT NULL REFERENCES problem(id) ON DELETE CASCADE,
    topic_id INT NOT NULL REFERENCES topic(id) ON DELETE CASCADE
);

INSERT INTO topic (topics)
VALUES ('golang'), ('rust'), ('js'), ('C');
