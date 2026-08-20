CREATE TABLE IF NOT EXISTS item_payloads (
    item_id INT NOT NULL REFERENCES items(id),
    ciphertext TEXT NOT NULL
);