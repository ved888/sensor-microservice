-- Revert password column length (assuming previous was VARCHAR(255))
ALTER TABLE users
    MODIFY COLUMN password VARCHAR(255) NOT NULL,
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS last_login;
