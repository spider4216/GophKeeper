CREATE TABLE IF NOT EXISTS sync_changes (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL, -- тут не будет реляций
    user_id INT NOT NULL,
    revision INT NOT NULL GENERATED ALWAYS AS IDENTITY, -- глобальный последовательный идентификатор изменения на сервере.
    operation VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);