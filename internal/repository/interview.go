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
		SELECT questions.public_id, questions.name, videos.public_id AS video_public_id, videos.path
		FROM questions
		JOIN positions ON questions.position_id = positions.id
		JOIN user_interviews ON user_interviews.position_id = positions.id
		JOIN interviews ON interviews.id = user_interviews.interview_id
		LEFT JOIN videos ON videos.interviews_public_id = interviews.public_id
		JOIN questions ON videos.question_public_id = questions.public_id
		WHERE interviews.public_id = $1;
	`

	result := models.InterviewResults{}
	rows, err := r.db.Query(ctx, query, publicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrInterviewNotFound
		}
		r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		question := models.Question{}
		var videoPublicId, videoPath string
		err := rows.Scan(&question.PublicID, &question.Name, &videoPublicId, &videoPath)
		if err != nil {
			r.logger.Errorf("Error occurred while scanning rows: %v", err)
			return nil, err
		}
		result.Result.Questions = append(result.Result.Questions, models.QuestionResult{
			Question:       question.Name,
			PublicID:       question.PublicID,
			VideoLink:      videoPath,
			VideoPublicID:  videoPublicId,
			EmotionResults: make([]models.EmotionResult, 0),
		})

	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating rows: %v", err)
		return nil, err
	}

	query = `SELECT c.public_id from candidates AS c
	JOIN user_interviews ui ON ui.candidate_id = c.id
	JOIN interviews i ON i.id = ui.interview_id
	WHERE i.public_id = $1 
	GROUP BY c.public_id`
	err = r.db.QueryRow(ctx, query, publicID).Scan(&result.CandidatePublicID)
	if err != nil {
		r.logger.Errorf("Error occurred while getting candidate public id: %v", err)
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

func (r *interviewRepository) AddVideoToQuestion(questionPublicID string, video string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	interviewPublicID := ""
	query := `SELECT i.public_id FROM interviews as i
	JOIN user_interviews ui ON ui.interview_id = i.id
	JOIN positions p ON p.id = ui.position_id
	JOIN questions q ON q.position_id = p.id
	WHERE q.public_id = $1
	`
	err := r.db.QueryRow(ctx, query, questionPublicID).Scan(&interviewPublicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrQuestionNotFound
		}
		r.logger.Errorf("Error occurred while retrieving add video to question: %v", err)
		return err
	}

	query = `
	INSERT INTO videos (interviews_public_id, question_public_id, path)
	VALUES ($1, $2, $3)
	RETURNING id;
`

	_, err = r.db.Exec(ctx, query, interviewPublicID, questionPublicID, video)
	if err != nil {
		r.logger.Errorf("Error occurred while adding video to question: %v", err)
		return err
	}

	return nil
}
