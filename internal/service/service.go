package service

import (
	"github.com/Zhiyenbek/sp-interview-main-service/config"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/repository"
	"go.uber.org/zap"
)

type InterviewsService interface {
	CreateInterviewResult(publicID string) (*models.InterviewResults, error)
}
type Service struct {
	InterviewsService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{}
}
