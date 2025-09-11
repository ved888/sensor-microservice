-- Revert password column length (assuming previous was VARCHAR(255))
ALTER TABLE users
    MODIFY COLUMN password VARCHAR(255) NOT NULL,
    DROP COLUMN first_name,
    DROP COLUMN last_name,
    DROP COLUMN last_login;
