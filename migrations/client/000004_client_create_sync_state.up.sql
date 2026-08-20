CREATE TABLE IF NOT EXISTS sync_state (
    last_server_revision INT NOT NULL,
    user_id INT NOT NULL
);