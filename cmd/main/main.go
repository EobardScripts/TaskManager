package main

import (
	"TaskManager/pkg/handlersService"
	"TaskManager/pkg/logger"
	"TaskManager/pkg/migrator"
	"TaskManager/pkg/storage"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	connStr := "postgres://test:test@192.168.1.102:3050/task_manager"
	storage, err := storage.New(connStr)
	if err != nil {
		logger.Error("Storage error: %s", err.Error())
	}

	handlerService := handlersService.New(storage)
	wg.Add(1)
	go handlerService.PreloadRoutes()

	//Миграции
	{
		err = migrator.Migration(storage)

		if err != nil {
			logger.Error("Migration error: %s", err.Error())
		}
	}

	wg.Wait()
}
