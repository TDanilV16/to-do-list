package controller

import (
	"fmt"
	"time"
	"to-do-list/internal/tasks"
)

type Storage interface {
	Create(task tasks.Task) error
	List(orderByAscending bool) (tasks.TaskList, error)
	ReadAllTasks() (tasks.TaskList, error)
	WriteToFile(list tasks.TaskList) error
}

type Controller struct {
	storage Storage
}

func (c *Controller) Delete(name string) error {
	list, err := c.storage.ReadAllTasks()
	if err != nil {
		return err
	}

	newList := make(tasks.TaskList, 0)

	for _, task := range list {
		if task.Name == name {
			continue
		}
		newList = append(newList, task)
	}

	return c.storage.WriteToFile(newList)
}

func (c *Controller) List(orderByAscending bool) error {
	list, err := c.storage.List(orderByAscending)
	if err != nil {
		return err
	}

	for _, task := range list {
		fmt.Printf(
			"status: %s, name: %s, description: %s, deadline: %s\n",
			task.Status,
			task.Name,
			task.Description,
			task.Deadline.Format("2006-01-02"))
	}

	return nil
}

func (c *Controller) Create(name, description, deadline string) error {
	t, err := time.Parse("2006-01-02", deadline)
	if err != nil {
		return err
	}

	return c.storage.Create(tasks.Task{
		Name:        name,
		Description: description,
		Deadline:    t,
		Status:      0,
	})
}

func NewController(storage Storage) *Controller {
	return &Controller{storage: storage}
}
