package storage

import (
	"TaskManager/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
)

var connStr = "postgres://test:test@192.168.1.102:3050/task_manager"

func newConnet() *pgxpool.Pool {
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		logger.Error("%s", err)
		return nil
	}

	return db
}

func TestStorage_Tasks(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		taskID   int
		authorID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TaskID: 1, AuthorID: 0",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				taskID:   1,
				authorID: 0,
			},
			wantErr: false,
		},
		{
			name: "TaskID: 14, AuthorID: 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				taskID:   14,
				authorID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.Tasks(tt.args.taskID, tt.args.authorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Log(fmt.Sprintf("Tasks() got = %+v", got))
		})
	}
}

func TestStorage_AllLabels(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Все метки",
			fields:  fields{DB: newConnet()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.AllLabels()
			if (err != nil) != tt.wantErr {
				t.Errorf("AllLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("AllLabels() got = %+v", got))
		})
	}
}

func TestStorage_AllTasks(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Все задачи",
			fields: fields{
				DB: newConnet(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.AllTasks()
			if (err != nil) != tt.wantErr {
				t.Errorf("AllTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("AllTasks() got = %+v", got))
		})
	}
}

func TestStorage_AllUsers(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Все пользователи",
			fields: fields{
				DB: newConnet(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.AllUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("AllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("AllUsers() got = %+v", got))
		})
	}
}

func TestStorage_DeleteLabel(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Удаление метки с id = 3",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 3,
			},
			wantErr: false,
		},
		{
			name: "Удаление метки с id = 4",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.DeleteLabel(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("DeleteLabel() got = %+v", got))
		})
	}
}

func TestStorage_DeleteTask(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Удаление задачи с id 14",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 14,
			},
			wantErr: false,
		},
		{
			name: "Удаление задачи с id 54",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 15,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.DeleteTask(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("DeleteTask() got = %+v", got))
		})
	}
}

func TestStorage_DeleteUser(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Удаление пользователя с id 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "Удаление пользователя с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.DeleteUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("DeleteUser() got = %+v", got))
		})
	}
}

func TestStorage_LabelById(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Поиск метки по id 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "Поиск метки по id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.LabelById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LabelById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("LabelById() got = %+v", got))
		})
	}
}

func TestStorage_NewLabel(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		label *Label
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Создание новой метки1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				label: &Label{
					Name: "Диагностика",
				},
			},
			wantErr: false,
		},
		{
			name: "Создание новой метки2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				label: &Label{
					Name: "Тестирование",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.NewLabel(tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("NewLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("NewLabel() got = %+v", tt.args.label))
		})
	}
}

func TestStorage_NewTask(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		t *Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Создание задачи1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				t: &Task{
					Title:   "Ремонт ПК",
					Content: "Проверить пк и отремонтировать",
				},
			},
			wantErr: false,
		},
		{
			name: "Создание задачи2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				t: &Task{
					Title:   "Диагностика ПК",
					Content: "Диагностика, по требованию клиента ремонт.",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.NewTask(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("NewTask() got = %+v", tt.args.t))
		})
	}
}

func TestStorage_NewTasks(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		tasks []*Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Создание задач из массива1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				tasks: []*Task{
					{Title: "Тест задачи1", Content: "Контент тестовой задачи1"},
					{Title: "Тест задачи2", Content: "Контент тестовой задачи2"},
					{Title: "Тест задачи3", Content: "Контент тестовой задачи3"},
				},
			},
			wantErr: false,
		},
		{
			name: "Создание задач из массива1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				tasks: []*Task{
					{Title: "Тест задачи4", Content: "Контент тестовой задачи4"},
					{Title: "Тест задачи5", Content: "Контент тестовой задачи5"},
					{Title: "Тест задачи6", Content: "Контент тестовой задачи6"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.NewTasks(tt.args.tasks); (err != nil) != tt.wantErr {
				t.Errorf("NewTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("NewTasks() got = %+v", tt.args.tasks))
		})
	}
}

func TestStorage_NewUser(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		user *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Создание пользователя1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				user: &User{
					Name: "Tester1",
				},
			},
			wantErr: false,
		},
		{
			name: "Создание пользователя2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				user: &User{
					Name: "Tester2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.NewUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("NewUser() got = %+v", tt.args.user))
		})
	}
}

func TestStorage_TaskById(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		taskID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Поиск задачи по id 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				taskID: 1,
			},
			wantErr: false,
		},
		{
			name: "Поиск задачи по id 4",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				taskID: 4,
			},
			wantErr: false,
		},
		{
			name: "Поиск задачи по id 13",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				taskID: 13,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.TaskById(tt.args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("TaskById() got = %+v", got))
		})
	}
}

func TestStorage_TasksByAuthor(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		authorID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Поиск задачи по автору с ID 0",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				authorID: 0,
			},
			wantErr: false,
		},
		{
			name: "Поиск задачи по автору с ID 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				authorID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.TasksByAuthor(tt.args.authorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TasksByAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("TasksByAuthor() got = %+v", got))
		})
	}
}

func TestStorage_TasksByLabel(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		labelID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Поиск задачи по метке с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				labelID: 2,
			},
			wantErr: false,
		},
		{
			name: "Поиск задачи по метке с id 3",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				labelID: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.TasksByLabel(tt.args.labelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TasksByLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("TasksByLabel() got = %+v", got))
		})
	}
}

func TestStorage_UpdateLabel(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		l *Label
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Обновлении метки с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				l: &Label{
					ID:   2,
					Name: "Чистка от пыли",
				},
			},
			wantErr: false,
		},
		{
			name: "Обновлении метки с id 3",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				l: &Label{
					ID:   3,
					Name: "Апгрейд",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.UpdateLabel(tt.args.l); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("UpdateLabel() got = %+v", tt.args.l))
		})
	}
}

func TestStorage_UpdateTask(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		t *Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Обновлении задачи с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				t: &Task{
					ID:      2,
					Title:   "Это обновленный title2",
					Content: "Это обновленный content2",
				},
			},
			wantErr: false,
		},
		{
			name: "Обновлении задачи с id 3",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				t: &Task{
					ID:      3,
					Title:   "Это обновленный title3",
					Content: "Это обновленный content3",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.UpdateTask(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("UpdateTask() got = %+v", tt.args.t))
		})
	}
}

func TestStorage_UpdateUser(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Обновление пользователя с id 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				u: &User{
					ID:   1,
					Name: "Иван (Обновлено)",
				},
			},
			wantErr: false,
		},
		{
			name: "Обновление пользователя с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				u: &User{
					ID:   2,
					Name: "Василий (Обновлено)",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			if err := s.UpdateUser(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(fmt.Sprintf("UpdateUser() got = %+v", tt.args.u))
		})
	}
}

func TestStorage_UserById(t *testing.T) {
	type fields struct {
		DB *pgxpool.Pool
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Получение пользователя с id 1",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "Получение пользователя с id 2",
			fields: fields{
				DB: newConnet(),
			},
			args: args{
				id: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				DB: tt.fields.DB,
			}
			got, err := s.UserById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(fmt.Sprintf("UserById() got = %+v", got))
		})
	}
}
