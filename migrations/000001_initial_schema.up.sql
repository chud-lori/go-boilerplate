CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    email VARCHAR(255) unique NOT NULL,
    password VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE posts (
    id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL, -- Added NOT NULL as it's typically required
    body TEXT,
    author_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- New updated_at column
    PRIMARY KEY (id),
    FOREIGN KEY (author_id) REFERENCES users (id)
);

-- Create a function to update the 'updated_at' timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger that calls the function before each UPDATE on the 'posts' table
CREATE TRIGGER set_posts_updated_at
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
