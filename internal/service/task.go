package service

import (
	"context"
	"todo-api-go/pkg/database"
	"todo-api-go/prisma/db"
)

type TaskService struct {
}

func (s *TaskService) GetAllTodos() ([]db.TaskModel, error) {
	tasks, err := database.PrismaClient.Task.FindMany().OrderBy().Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) Create(title, description, status string) (db.TaskModel, error) {
	todo, err := database.PrismaClient.Task.CreateOne(
		db.Task.Title.Set(title),
		db.Task.Description.Set(description),
		db.Task.Status.Set(db.Status(status)),
	).Exec(context.Background())

	if err != nil {
		return db.TaskModel{}, err
	}

	return *todo, nil
}

func (s *TaskService) GetById(id string) (db.TaskModel, error) {
	todo, err := database.PrismaClient.Task.FindUnique(
		db.Task.ID.Equals(id),
	).Exec(context.Background())

	if err != nil {
		return db.TaskModel{}, err
	}

	return *todo, nil
}
