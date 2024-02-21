package migrator

import (
	"TaskManager/pkg/storage"
	"context"
)

func Migration(storage *storage.Storage) error {
	ctx := context.Background()
	tx, err := storage.DB.Begin(ctx)
	if err != nil {
		return err
	}
	//1ая строка запроса для тестов: DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;
	_, err = tx.Exec(context.Background(), `
		DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

		CREATE TABLE IF NOT EXISTS users (
    	id SERIAL PRIMARY KEY,
    	name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS labels (
    		id SERIAL PRIMARY KEY,
    		name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS tasks (
    		id SERIAL PRIMARY KEY,
    		opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
    		closed BIGINT DEFAULT 0,
    		author_id INTEGER REFERENCES users(id) DEFAULT 1,
    		assigned_id INTEGER REFERENCES users(id) DEFAULT 1,
    		title TEXT,
    		content TEXT
		);

		CREATE TABLE IF NOT EXISTS tasks_labels (
    		task_id INTEGER REFERENCES tasks(id),
    		label_id INTEGER REFERENCES labels(id)
		);

		INSERT INTO users (name) VALUES ('default');
	`)

	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
