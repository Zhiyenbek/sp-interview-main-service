package handler

import (
	"net/http"

	"github.com/Zhiyenbek/sp-interview-main-service/internal/models"
	"github.com/gin-gonic/gin"
)

// func (h *handler) UploadVideo(c *gin.Context) {
// 	interviewID := c.Param("id")

// 	// Parse the JSON request body
// 	req := models.VideoRequest{}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	filePath := fmt.Sprintf("/%s/interview_%s_question_%d.mp4", h.cfg.Video.Path, interviewID, req.QuestionNumber)
// 	err := os.WriteFile(filePath, req.VideoFile, 0644)
// 	if err != nil {
// 		h.logger.Errorf("could not save video file %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video file"})
// 		return
// 	}
// 	//TODO add to save it db
// 	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
// }

func (h *handler) CreateInterviewResult(c *gin.Context) {
	interviewID := c.Param("id")
	res, err := h.service.CreateInterviewResult(interviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusCreated, sendResponse(0, res, nil))
}
