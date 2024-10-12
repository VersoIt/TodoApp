package handler

import (
	"TodoApp/internal/model"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary createList
// @Security ApiKeyAuth
// @Tags list
// @Description create list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body model.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.TodoList
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.CreateList(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// @Summary getAllLists
// @Security ApiKeyAuth
// @Tags list
// @Description get all lists
// @ID get-all-lists
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"data": lists})
}

// @Summary getListById
// @Security ApiKeyAuth
// @Tags list
// @Description get list by ids
// @ID get-list-by-id
// @Produce json
// @Param id path int true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "list not found")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary updateList
// @Security ApiKeyAuth
// @Tags list
// @Description update lists
// @ID update-list
// @Accept json
// @Produce json
// @Param id path int true "list id"
// @Param input body model.UpdateListInput true "list info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.UpdateListInput
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoList.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "success"})
}

// @Summary deleteList
// @Security ApiKeyAuth
// @Tags list
// @Description delete list
// @ID delete-list
// @Produce json
// @Param id path int true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "success",
	})
}
