CREATE TABLE IF NOT EXISTS pending_changes (
    item_id CHAR(36) NOT NULL,
    operation VARCHAR(50) NOT NULL, -- CREATE | UPDATE | DELETE
    user_id INT NOT NULL
);
