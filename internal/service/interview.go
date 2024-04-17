package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

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
type QuestionReq struct {
	Question  string `json:"question"`
	PublicID  string `json:"public_id"`
	VideoLink string `json:"video_link"`
}

type Request struct {
	Questions []QuestionReq `json:"questions"`
}
type Result struct {
	Result models.Result `json:"result"`
}

func NewInterviewsService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *interviewsService {
	return &interviewsService{
		interviewRepo: repo.InterviewRepository,
		cfg:           cfg,
		logger:        logger,
	}
}

func (s *interviewsService) AddVideoToQuestion(questionPublicID, interviewPublicID, video string) error {
	return s.interviewRepo.AddVideoToQuestion(questionPublicID, interviewPublicID, video)
}

func (s *interviewsService) CreateInterviewResult(publicID string) (*models.InterviewResults, error) {
	interview, err := s.interviewRepo.GetInterviewByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	req := Request{
		Questions: make([]QuestionReq, 0),
	}

	for _, q := range interview.Result.Questions {
		req.Questions = append(req.Questions, QuestionReq{
			PublicID:  q.PublicID,
			Question:  q.Question,
			VideoLink: q.VideoLink,
		})
	}
	res, err := sendDataToAPI(req, s.cfg.Video.Url)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	if res != nil {
		interview.Result = res.Result
		if len(res.Result.Questions) != 0 {
			interview.Result.Score = res.Result.Score / len(res.Result.Questions)
		}
	}

	interview.RawResult, err = json.Marshal(res.Result)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	interview.PublicID = publicID

	err = s.interviewRepo.PutInterview(interview)
	if err != nil {
		return nil, err
	}

	return interview, nil
}

func (s *interviewsService) GetInterviewByPublicID(publicID string) (*models.InterviewResults, error) {
	return s.interviewRepo.GetInterview(publicID)
}

func (s *interviewsService) GetAllInterviews() ([]*models.InterviewResults, error) {
	return s.interviewRepo.GetAllInterviews()
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 600*time.Second)
}

func sendDataToAPI(data Request, url string) (*Result, error) {
	transport := http.Transport{
		Dial: dialTimeout,
	}
	// Convert the InterviewResults struct to JSON
	client := &http.Client{
		Transport: &transport, // Set the timeout duration here
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
	}
	// Create a request body with the JSON data
	body := bytes.NewReader(jsonData)

	// Send a POST request to the API endpoint
	resp, err := client.Post(url+"/process_interview", "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to API: %v", err)
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	// Check the response status code
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnprocessableEntity {
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", respBody)
			}
			return nil, fmt.Errorf("API request failed with status code %d. %v", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println(string(respBody))
	var responseData Result
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &responseData, nil
}
