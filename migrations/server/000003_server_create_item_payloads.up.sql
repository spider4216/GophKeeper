CREATE TABLE IF NOT EXISTS item_payloads (
    item_id CHAR(36) NOT NULL REFERENCES items(id),
    ciphertext TEXT NOT NULL
);