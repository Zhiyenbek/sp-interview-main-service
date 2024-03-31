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

type QuestionResult struct {
	Question       string          `json:"question"`
	PublicID       string          `json:"public_id"`
	QuestionType   string          `json:"question_type"`
	Evaluation     string          `json:"evaluation"`
	Score          int             `json:"score"`
	Answer         string          `json:"answer"`
	Emotion        string          `json:"emotion"`
	VideoLink      string          `json:"video_link"`
	VideoPublicID  string          `json:"video_public_id"`
	EmotionResults []EmotionResult `json:"emotion_results"`
}

type EmotionResult struct {
	Emotion   string  `json:"emotion"`
	ExactTime float64 `json:"exact_time"`
	Duration  float64 `json:"duration"`
}

type Result struct {
	Questions []QuestionResult `json:"questions"`
	Score     int              `json:"score"`
}

type Question struct {
	Name     string
	PublicID string
}
