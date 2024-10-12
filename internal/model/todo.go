package model

import "errors"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UserList struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	ListId int `json:"list_id"`
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListItem struct {
	Id     int `json:"id"`
	ListId int `json:"list_id"`
	ItemId int `json:"item_id"`
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (u UpdateListInput) Validate() error {
	if u.Title == nil && u.Description == nil {
		return errors.New("either title or description must be set")
	}

	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (u UpdateItemInput) Validate() error {
	if u.Title == nil && u.Description == nil && u.Done == nil {
		return errors.New("either title, description or done must be set")
	}

	return nil
}
