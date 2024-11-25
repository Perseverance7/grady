package models

type CreateGroupReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	CreatedBy   int64
}
