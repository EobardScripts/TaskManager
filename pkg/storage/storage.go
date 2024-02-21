package storage

import (
	"TaskManager/pkg/logger"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	DB *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		DB: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Пользователь
type User struct {
	ID   int
	Name string
}

// Метки
type Label struct {
	ID   int
	Name string
}

// -------------------Метки-------------------------

// NewLabel - создание новой метки, возвращает все поля новой метки
func (s *Storage) NewLabel(label *Label) error {
	var id int
	err := s.DB.QueryRow(context.Background(), `
		INSERT INTO labels (name)
		VALUES ($1) RETURNING id;
		`,
		label.Name,
	).Scan(&id)

	if err != nil {
		return err
	}

	thisLabel, err := s.LabelById(id)
	if err != nil {
		return err
	}

	*label = *thisLabel

	return nil
}

// LabelById - находит и возвращает метку по id
func (s *Storage) LabelById(id int) (*Label, error) {
	label := &Label{}
	err := s.DB.QueryRow(context.Background(), `
		SELECT id, name 
		FROM labels
		WHERE id = $1;
		`,
		id,
	).Scan(&label.ID, &label.Name)

	if err != nil {
		return label, err
	}

	return label, nil
}

// AllLabels - Возвращает все метки
func (s *Storage) AllLabels() ([]Label, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			id,
			name
		FROM labels
		ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var labels []Label
	for rows.Next() {
		var l Label
		err = rows.Scan(
			&l.ID,
			&l.Name)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		labels = append(labels, l)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return labels, rows.Err()
}

// UpdateLabel - обновляет метку и возвращает уже обновленную модель
func (s *Storage) UpdateLabel(l *Label) error {
	_, err := s.DB.Exec(context.Background(), `
		UPDATE labels
		SET name = $1
		WHERE
			(id = $2);`,
		l.Name,
		l.ID,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении метки: %s", err.Error())
		return err
	}

	thisLabel, err := s.LabelById(l.ID)
	if err != nil {
		return err
	}

	*l = *thisLabel

	return err
}

// DeleteLabel - удаляет метку по ее ID и возвращает удаленную запись
func (s *Storage) DeleteLabel(id int) (*Label, error) {
	thisLabel, err := s.LabelById(id)
	if err != nil {
		return thisLabel, err
	}
	_, err = s.DB.Exec(context.Background(), `
		DELETE FROM labels
		WHERE
			(id = $1);`,
		id,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении метки: %s", err.Error())
		return thisLabel, err
	}

	return thisLabel, nil
}

//-------------------Пользователи-------------------------

// NewUser - создание нового пользователя, возвращает все поля нового пользователя
func (s *Storage) NewUser(user *User) error {
	var id int
	err := s.DB.QueryRow(context.Background(), `
		INSERT INTO users (name)
		VALUES ($1) RETURNING id;
		`,
		user.Name,
	).Scan(&id)

	if err != nil {
		return err
	}

	thisUser, err := s.UserById(id)
	if err != nil {
		return err
	}

	*user = *thisUser

	return nil
}

// UserById - находит и возвращает пользователя по id
func (s *Storage) UserById(id int) (*User, error) {
	user := &User{}
	err := s.DB.QueryRow(context.Background(), `
		SELECT id, name 
		FROM users
		WHERE id = $1;
		`,
		id,
	).Scan(&user.ID, &user.Name)

	if err != nil {
		return user, err
	}

	return user, nil
}

// AllUsers - Возвращает всех пользователей
func (s *Storage) AllUsers() ([]User, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			id,
			name
		FROM users
		ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.Name)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		users = append(users, u)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return users, rows.Err()
}

// UpdateUser - обновляет пользователя и возвращает уже обновленную модель
func (s *Storage) UpdateUser(u *User) error {
	_, err := s.DB.Exec(context.Background(), `
		UPDATE users
		SET name = $1
		WHERE
			(id = $2);`,
		u.Name,
		u.ID,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении пользователя: %s", err.Error())
		return err
	}

	thisUser, err := s.UserById(u.ID)
	if err != nil {
		return err
	}

	*u = *thisUser

	return err
}

// DeleteUser - удаляет пользователя по его ID и возвращает удаленную запись
func (s *Storage) DeleteUser(id int) (*User, error) {
	thisUser, err := s.UserById(id)
	if err != nil {
		return thisUser, err
	}
	_, err = s.DB.Exec(context.Background(), `
		DELETE FROM users
		WHERE
			(id = $1);`,
		id,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении пользователя: %s", err.Error())
		return thisUser, err
	}

	return thisUser, nil
}

//-------------------Задачи-------------------------

// TaskById - возвращает задачу по ее id
func (s *Storage) TaskById(taskID int) (*Task, error) {
	t := Task{}
	row := s.DB.QueryRow(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			id = $1
		ORDER BY id;
	`, taskID)

	err := row.Scan(
		&t.ID,
		&t.Opened,
		&t.Closed,
		&t.AuthorID,
		&t.AssignedID,
		&t.Title,
		&t.Content,
	)
	if err != nil {
		return &t, err
	}

	return &t, nil
}

// AllTasks - Возвращает все задачи
func (s *Storage) AllTasks() ([]Task, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// TasksByLabel - возвращает список задач по ID метки
func (s *Storage) TasksByLabel(labelID int) ([]Task, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			t.id,
			t.opened,
			t.closed,
			t.author_id,
			t.assigned_id,
			t.title,
			t.content
		FROM tasks as t
		INNER JOIN tasks_labels as tl
    	ON (tl.label_id = $1) AND (t.id = tl.task_id);`,
		labelID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// TasksByAuthor - возвращает список задач по ID автора
func (s *Storage) TasksByAuthor(authorID int) ([]Task, error) {
	rows, err := s.DB.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
		author_id = $1
		ORDER BY id;
	`,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask - создаёт новую задачу и возвращает все поля в t *Task.
func (s *Storage) NewTask(t *Task) error {
	var id int
	err := s.DB.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)

	thisTask, err := s.TaskById(id)
	if err != nil {
		return err
	}

	*t = *thisTask

	return err
}

// NewTasks - создаёт массив задач и возвращает все поля в t []*Task.
func (s *Storage) NewTasks(tasks []*Task) error {
	ctx := context.Background()
	tx, err := s.DB.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Prepare(ctx, "my-insert", `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
	`)

	if err != nil {
		logger.Error("Ошибка при подготовке плана: %s", err.Error())
		return err
	}

	for _, task := range tasks {
		row := tx.QueryRow(ctx, "my-insert", task.Title, task.Content)
		err := row.Scan(&task.ID)
		if err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

// UpdateTask - обновляет задачу и возвращает уже обновленную модель
func (s *Storage) UpdateTask(t *Task) error {
	_, err := s.DB.Exec(context.Background(), `
		UPDATE tasks
		SET (title, content, closed) = ($1, $2, $3)
		WHERE
			(id = $4);`,
		t.Title,
		t.Content,
		t.Closed,
		t.ID,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении задачи: %s", err.Error())
		return err
	}

	thisTask, err := s.TaskById(t.ID)
	if err != nil {
		return err
	}

	*t = *thisTask

	return err
}

// DeleteTask - удаляет задачу по ее ID и возвращает удаленную запись
func (s *Storage) DeleteTask(id int) (*Task, error) {
	thisTask, err := s.TaskById(id)
	if err != nil {
		return thisTask, err
	}
	_, err = s.DB.Exec(context.Background(), `
		DELETE FROM tasks
		WHERE
			(id = $1);`,
		id,
	)

	if err != nil {
		logger.Error("Ошибка при обновлении задачи: %s", err.Error())
		return thisTask, err
	}

	return thisTask, nil
}
