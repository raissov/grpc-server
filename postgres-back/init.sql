CREATE TABLE IF NOT EXISTS message (
    id serial PRIMARY KEY,
    message text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
)