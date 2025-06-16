CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    email VARCHAR(255) unique NOT NULL,
    passcode VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
