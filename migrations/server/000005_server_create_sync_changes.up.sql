CREATE TABLE IF NOT EXISTS sync_changes (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES items(id),
    revision INT NOT NULL,
    user_id INT NOT NULL REFERENCES users(id),
    operation VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);

/*
todo индексы когда пойму по каким полям чаще делаю выборку
*/