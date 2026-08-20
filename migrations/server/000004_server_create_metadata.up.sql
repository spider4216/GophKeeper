CREATE TABLE IF NOT EXISTS metadata (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES items(id),
    key VARCHAR(50) NOT NULL,
    value VARCHAR(255) NOT NULL
);

/*
todo индексы когда пойму по каким полям чаще делаю выборку
*/