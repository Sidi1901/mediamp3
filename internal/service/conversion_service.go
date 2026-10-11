package service 

import(
	"errors"
	"github.com/Sidi1901/mediamp3/internal/model"
)

var ErrConversionNotImplemented = errors.New("conversion functionality is not implemented yet")

type ConversionService interface {
    ConvertFormat(request model.ConversionRequest) (model.ConversionResponse, error)  
}

type conversionService struct {}

func NewConversionService() ConversionService {
	return &conversionService{}
}

func (s *conversionService) ConvertFormat(request model.ConversionRequest) (model.ConversionResponse, error) {
	return model.ConversionResponse{}, ErrConversionNotImplemented
}