CREATE TABLE IF NOT EXISTS pending_changes (
    item_id INT NOT NULL REFERENCES items(id),
    operation VARCHAR(50) NOT NULL -- CREATE | UPDATE | DELETE
);

/*
todo индексы когда пойму по каким полям чаще делаю выборку
*/