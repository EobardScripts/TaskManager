package storage

import (
	"TaskManager/pkg/logger"
	"testing"
)

var connStr = "postgres://test:test@192.168.1.102:3050/task_manager"

var s *Storage

func Init() {
	s, err := New(connStr)
	if err != nil {
		logger.Error("%s", err)
	}
	_ = s
}

func TestStorage_Tasks(t *testing.T) {
	Init()
	tasks, err := s.Tasks(1, 0)
	if err != nil {
		logger.Error("Ошибка при получении списка задач: %s", tasks)
		t.Fatal(err)
	}
	t.Log(tasks)
	tasks, err = s.Tasks(14, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tasks)
}
