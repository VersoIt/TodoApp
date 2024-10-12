package handler

import (
	"TodoApp/internal/model"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary createItem
// @Security ApiKeyAuth
// @Tags item
// @Description creates an item
// @ID create-item
// @Accept json
// @Produce json
// @Param input body model.TodoItem true "item info"
// @Param id path int true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input model.TodoItem
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary getAllItems
// @Security ApiKeyAuth
// @Tags item
// @Description get all items
// @ID get-all-items
// @Produce json
// @Param id path int true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary getItemById
// @Security ApiKeyAuth
// @Tags item
// @Description get item by id
// @ID get-item-by-id
// @Produce json
// @Param id path int true "item id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "item not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary updateItem
// @Security ApiKeyAuth
// @Tags item
// @Description update item by id
// @ID update-item-by-id
// @Accept json
// @Produce json
// @Param id path int true "list id"
// @Param input body model.UpdateItemInput true "item info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var updateItemInput model.UpdateItemInput
	if err = c.BindJSON(&updateItemInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoItem.Update(userId, itemId, updateItemInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "item not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}

// @Summary deleteItem
// @Security ApiKeyAuth
// @Tags item
// @Description delete item by id
// @ID update-item-by-id
// @Produce json
// @Param id path int true "item id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "item not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}
