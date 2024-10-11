CREATE TABLE IF NOT EXISTS metadata(
    id serial PRIMARY KEY,
    title VARCHAR (300),
    explanation TEXT,
    "date" DATE UNIQUE,
    copyright VARCHAR (300),
    filename VARCHAR (14)
);

CREATE INDEX IF NOT EXISTS idx_metadata_date ON metadata("date");
