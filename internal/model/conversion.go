package model 

type ConversionRequest struct {
	URL     string  `json:"url" binding:"required"`
	Format  string  `json:"format" binding:"requires,oneof=mp3 mp4"`
	Quality string  `json:"quality" binding:"required"`
}

type ConversionStatus string 

const (
	ConversionStatusQueued     ConversionStatus = "queued"
	ConversionStatusProcessing ConversionStatus = "processing"
	ConversionStatusCompleted  ConversionStatus = "completed"
	ConversionStatusFailed     ConversionStatus = "failed"
)


type ConversionJob struct {
	ID      string            `json:"id"`
	URL     string            `json:"url"`
	Format  string            `json:"format"`
	Quality string            `json:"quality"`
	Status  ConversionStatus  `json:"status"`
}


type ConversionResponse struct {
	ID     string            `json:"id"`
	status ConversionStatus  `json:"status"` 
}