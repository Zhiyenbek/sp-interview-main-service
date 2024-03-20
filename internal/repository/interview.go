package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Zhiyenbek/sp-interview-main-service/config"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type interviewRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewInterviewRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) InterviewRepository {
	return &interviewRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *interviewRepository) GetInterviewByPublicID(publicID string) (*models.InterviewResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		SELECT public_id, results
		FROM interviews
		WHERE public_id = $1;
	`

	result := models.InterviewResults{}
	err := r.db.QueryRow(ctx, query, publicID).Scan(
		&result.PublicID,
		&result.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrInterviewNotFound
		}
		r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
		return nil, err
	}

	return &result, nil
}
func (r *interviewRepository) PutInterview(interview *models.InterviewResults) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		UPDATE interviews
		SET results = $1
		WHERE public_id = $2;
	`

	// Convert the interview results to JSON
	jsonData, err := json.Marshal(interview.Result)
	if err != nil {
		r.logger.Errorf("Failed to marshal interview results to JSON: %v", err)
		return err
	}

	// Execute the update query
	_, err = r.db.Exec(ctx, query, jsonData, interview.PublicID)
	if err != nil {
		r.logger.Errorf("Error occurred while updating interview results: %v", err)
		return err
	}

	return nil
}
