CREATE TABLE IF NOT EXISTS item_chunks (
    item_id CHAR(36) NOT NULL,
    chunk_number INT NOT NULL,
    ciphertext TEXT NOT NULL
);