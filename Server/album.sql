


CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    artist VARCHAR(100) NOT NULL,
    year INT,
    user_id INT
);

ALTER TABLE albums
ADD COLUMN IF NOT EXISTS Cover_url VARCHAR(100);
DROP COLUMN IF EXISTS Url;


CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL
);


INSERT INTO albums (title, artist, year)
VALUES
  ('Thriller', 'Michael Jackson', 1982),
  ('Back in Black', 'AC/DC', 1980),
  ('The Dark Side of the Moon', 'Pink Floyd', 1973),
  ('Nevermind', 'Nirvana', 1991),
  ('21', 'Adele', 2011);

SELECT * FROM albums;
SELECT title, artist, year FROM albums;

SELECT * FROM users;

SELECT 1 FROM users WHERE email = 'paulke@mail.co'

INSERT INTO albums (title, artist, year) VALUES ('PDP', 'B.B.Jacques', 2022) RETURNING id

SELECT EXISTS(SELECT 1 FROM users WHERE id = 1)