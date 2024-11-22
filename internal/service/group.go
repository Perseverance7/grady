package service

import (
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
)

type GroupService struct {
	repo repository.Group
}

func NewGroupService(repo repository.Group) *GroupService {
	return &GroupService{
		repo: repo,
	}
}

func (g *GroupService) CreateGroup(group *models.CreateGroupReq) error {
	return g.repo.CreateGroup(group)
}
