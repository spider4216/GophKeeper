CREATE TABLE IF NOT EXISTS pending_changes (
    -- item_id INT NOT NULL REFERENCES items(id),
    item_id INT NOT NULL, -- Подумать над soft deleted в items, чтобы использовать reference
    operation VARCHAR(50) NOT NULL, -- CREATE | UPDATE | DELETE
    user_id INT NOT NULL
);
