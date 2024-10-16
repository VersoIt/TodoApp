package service

import (
	"TodoApp/internal/model"
	"TodoApp/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) CreateList(userId int, list model.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	lists, err := s.repo.GetAll(userId)
	if lists == nil {
		lists = make([]model.TodoList, 0)
	}
	return lists, err
}

func (s *TodoListService) GetById(userId, listId int) (model.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, updateRequest model.UpdateListInput) error {
	if err := updateRequest.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, updateRequest)
}
