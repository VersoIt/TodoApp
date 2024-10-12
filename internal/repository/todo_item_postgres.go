package repository

import (
	"TodoApp/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItemRepository struct {
	db *sqlx.DB
}

func NewTodoItemRepository(db *sqlx.DB) *TodoItemRepository {
	return &TodoItemRepository{db: db}
}

func (r *TodoItemRepository) Create(listId int, todoItem model.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemsQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	err = r.db.QueryRow(createItemsQuery, todoItem.Title, todoItem.Description).Scan(&itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemRepository) GetAll(userId, listId int) ([]model.TodoItem, error) {
	var items []model.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id
									INNER JOIN %s ul ON li.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemRepository) GetById(userId, itemId int) (model.TodoItem, error) {
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s il ON il.item_id = ti.id INNER JOIN %s ul ON ul.list_id = il.list_id WHERE ul.user_id = $1 AND ti.id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	var item model.TodoItem
	if err := r.db.Get(&item, query, userId, itemId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemRepository) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul 
       								WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemRepository) Update(userId, itemId int, updateItemInput model.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if updateItemInput.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *updateItemInput.Title)
		argId++
	}

	if updateItemInput.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *updateItemInput.Description)
		argId++
	}

	if updateItemInput.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *updateItemInput.Done)
		argId++
	}

	setValuesQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s il, %s ul WHERE il.item_id = ti.id AND il.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d", todoItemsTable, setValuesQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
