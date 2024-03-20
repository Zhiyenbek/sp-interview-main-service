package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zhiyenbek/sp-interview-main-service/config"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/Zhiyenbek/sp-interview-main-service/internal/repository"
	"go.uber.org/zap"
)

type interviewsService struct {
	cfg           *config.Configs
	logger        *zap.SugaredLogger
	interviewRepo repository.InterviewRepository
}

func NewInterviewsService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *interviewsService {
	return &interviewsService{
		interviewRepo: repo.InterviewRepository,
		cfg:           cfg,
		logger:        logger,
	}
}

func (s *interviewsService) CreateInterviewResult(publicID string) (*models.InterviewResults, error) {
	res, err := s.interviewRepo.GetInterviewByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.RawResult, &res.Result)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	res, err = sendDataToAPI(res, s.cfg.Video.Url)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	res.RawResult, err = json.Marshal(res.Result)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	err = s.interviewRepo.PutInterview(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *interviewsService) GetInterviewByPublicID(publicID string) (*models.InterviewResults, error) {
	return s.interviewRepo.GetInterviewByPublicID(publicID)
}

func sendDataToAPI(data *models.InterviewResults, url string) (*models.InterviewResults, error) {
	// Convert the InterviewResults struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	// Create a request body with the JSON data
	body := bytes.NewReader(jsonData)

	// Send a POST request to the API endpoint
	resp, err := http.Post(url+"/proccess_interview", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to API: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData models.InterviewResults
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &responseData, nil
}
