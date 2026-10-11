package service

import(
	"github.com/Sidi1901/mediamp3/internal/model"
)

type FormatService interface {
	GetFormats() model.FormatResponse
}

type formatService struct {}

func NewFormatService() FormatService {
	return &formatService{}
}

func (s *formatService) GetFormats() model.FormatResponse {
	
	modelResponse := model.FormatResponse{
		Formats:  []model.Format{
			{
				Type: "mp3",
				Qualities: []string{
					"128kbps",
					"192kbps",
					"320kbps",
				},
			},
			{
				Type: "mp4",
				Qualities: []string{
					"360p",
					"480p",
					"720p",
					"1080p",
				},
			},
		},
	}

	return modelResponse
}