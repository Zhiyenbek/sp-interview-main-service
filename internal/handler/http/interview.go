package handler

import (
	"errors"
	"net/http"

	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// func (h *handler) UploadVideo(c *gin.Context) {
// 	interviewID := c.Param("id")

// 	// Parse the JSON request body
// 	req := models.VideoRequest{}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

//		filePath := fmt.Sprintf("/%s/interview_%s_question_%d.mp4", h.cfg.Video.Path, interviewID, req.QuestionNumber)
//		err := os.WriteFile(filePath, req.VideoFile, 0644)
//		if err != nil {
//			h.logger.Errorf("could not save video file %v", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video file"})
//			return
//		}
//		//TODO add to save it db
//		c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
//	}
type Video struct {
	Video             string `json:"video" binding:"required"`
	InterviewPublicID string `json:"interview_public_id" binding:"required"`
}

func (h *handler) CreateInterviewResult(c *gin.Context) {
	interviewID := c.Param("id")
	res, err := h.service.CreateInterviewResult(interviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusCreated, sendResponse(0, res, nil))
}

func (h *handler) AddVideoToQuestion(c *gin.Context) {
	questionID := c.Param("id")
	req := &Video{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("Failed to parse request body when deleting skills from position: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	err := h.service.InterviewsService.AddVideoToQuestion(questionID, req.InterviewPublicID, req.Video)
	if err != nil {
		if errors.Is(err, models.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, sendResponse(-1, nil, models.ErrQuestionNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) GetInterviews(c *gin.Context) {
	res, err := h.service.InterviewsService.GetAllInterviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, res, nil))
}

func (h *handler) GetInterviewByPublicID(c *gin.Context) {
	publicID := c.Param("interview_public_id")

	res, err := h.service.InterviewsService.GetInterviewByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, sendResponse(0, res, nil))

}
