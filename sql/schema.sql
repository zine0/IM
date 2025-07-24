CREATE TABLE "user" (
    ID SERIAL PRIMARY KEY,
    username VARCHAR(255),
    password VARCHAR(255),
    avatar VARCHAR(255),
    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_user_username ON "user"(username);
