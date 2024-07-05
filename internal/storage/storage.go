package storage

import (
	"encoding/json"
	"io"
	"slices"
	"to-do-list/internal/tasks"
)

type File interface {
	io.ReadWriteCloser
	Truncate(size int64) error
	Seek(offset int64, whence int) (ret int64, err error)
}

type Storage struct {
	file File
}

func (s *Storage) List(orderByAscending bool) (tasks.TaskList, error) {
	list, err := s.ReadAllTasks()
	if err != nil {
		return nil, err
	}
	if orderByAscending {
		slices.SortStableFunc(list, tasks.OrderByDeadlineAscending)
	} else {
		slices.SortStableFunc(list, tasks.OrderByDeadlineDescending)

	}
	return list, nil
}

func NewStorage(file File) *Storage {
	return &Storage{file: file}
}

func (s *Storage) WriteToFile(tasksList tasks.TaskList) error {
	bytes, err := json.Marshal(tasksList)
	if err != nil {
		return err
	}

	err = s.file.Truncate(0)
	if err != nil {
		return err

	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = s.file.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Create(task tasks.Task) error {
	tasksList, err := s.ReadAllTasks()
	if err != nil {
		return err
	}

	tasksList = append(tasksList, &task)

	return s.WriteToFile(tasksList)
}

func (s *Storage) ReadAllTasks() (tasks.TaskList, error) {
	readAll, err := io.ReadAll(s.file)
	if err != nil {
		return nil, err
	}

	var list tasks.TaskList

	if len(readAll) == 0 {
		return list, nil
	}

	if err := json.Unmarshal(readAll, &list); err != nil {
		return nil, err
	}

	return list, nil
}
