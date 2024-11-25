package handler

import (
	"net/http"

	"github.com/Perseverance7/grady/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createGroup(c *gin.Context) {
	var groupReq models.CreateGroupReq

	user, err := getUserInfo(c)
	if err != nil {
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	groupReq.CreatedBy = user.ID

	if err := c.BindJSON(&groupReq); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.CreateGroup(&groupReq)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "group succesfully created"})

}

func (h *Handler) listGroups(c *gin.Context) {

}

func (h *Handler) getGroupDetails(c *gin.Context) {

}

func (h *Handler) joinMember(c *gin.Context) {

}

func (h *Handler) removeMember(c *gin.Context) {

}
