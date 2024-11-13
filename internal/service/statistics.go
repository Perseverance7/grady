package service

import (
	"github.com/Perseverance7/grady/internal/repository"
)

type StatisticsService struct {
	repo repository.Statistics
}

func NewStatisticsService(repo repository.Statistics) *StatisticsService {
	return &StatisticsService{
		repo: repo,
	}
}
