package repository

import (
	"github.com/Zhiyenbek/sp-interview-main-service/config"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type InterviewRepository interface {
	GetInterviewByPublicID(publicID string) (*models.InterviewResults, error)
	PutInterview(interview *models.InterviewResults) error
	AddVideoToQuestion(questionPublicID string, video string) error
}
type Repository struct {
	InterviewRepository
}

func New(db *pgxpool.Pool, cfg *config.Configs, log *zap.SugaredLogger) *Repository {
	return &Repository{
		InterviewRepository: NewInterviewRepository(db, cfg.DB, log),
	}
}
