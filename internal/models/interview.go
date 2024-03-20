package models

type VideoRequest struct {
	QuestionNumber int    `json:"questionNumber"`
	VideoFile      []byte `json:"videoFile"`
}

type InterviewResults struct {
	PublicID  string `json:"public_id"`
	Result    Result `json:"result"`
	RawResult []byte `json:"-"`
}

type Question struct {
	Question       string          `json:"question"`
	QuestionType   string          `json:"question_type"`
	Evaluation     string          `json:"evaluation"`
	Score          int             `json:"score"`
	VideoLink      string          `json:"video_link"`
	EmotionResults []EmotionResult `json:"emotion_results"`
}

type EmotionResult struct {
	Emotion   string  `json:"emotion"`
	ExactTime float64 `json:"exact_time"`
	Duration  float64 `json:"duration"`
}

type Result struct {
	Questions []Question `json:"questions"`
	Score     int        `json:"score"`
}
