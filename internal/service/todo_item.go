package service

import (
	"TodoApp/internal/model"
	"TodoApp/internal/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, todoItem model.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// list does not exist or does not belong to user
		return 0, err
	}

	return s.repo.Create(listId, todoItem)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]model.TodoItem, error) {
	items, err := s.repo.GetAll(userId, listId)
	if items == nil {
		items = make([]model.TodoItem, 0)
	}
	return items, err
}

func (s *TodoItemService) GetById(userId, itemId int) (model.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, updateItemInput model.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, updateItemInput)
}
