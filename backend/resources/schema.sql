CREATE TABLE users (
    id             uuid   PRIMARY KEY,
    name           text   NOT NULL,
    games_played   int    NOT NULL DEFAULT 0,
    score          int    NOT NULL DEFAULT 0,
    friends        uuid[] NOT NULL DEFAULT array[]::uuid[]
);