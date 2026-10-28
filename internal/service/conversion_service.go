package service 

import(
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"github.com/Sidi1901/mediamp3/internal/model"
)

var (
	ErrInvalidConversionRequest = errors.New(("invalid conversion request")
	ErrConversionNotImplemented = errors.New("conversion functionality is not implemented yet")
)

type ConversionService interface {
    ConvertFormat(ctx context.Context, request model.ConversionRequest) (model.ConversionResponse, error)  
}

type conversionService struct {
	cfg config.Config
}

func NewConversionService(cfg config.Config) ConversionService {
	return &conversionService{
		cfg : cfg,
	}
}

func (s *conversionService) ConvertFormat(ctx context.Context, request model.ConversionRequest) (model.ConversionResponse, error) {
	
	request.URL := strings.Trimpspace(request.URL)
	request.Format = strings.ToLower(strings.TrimSpace(request.Format))
	request.Quality = strings.ToLower(strings.TrimSpace(request.Quality))

	if err := validateConversionRequest(request); err != nil {
		return model.ConversionResponse{}, err
	}

	id, err := generateID()

	if err != nil {
		return model.ConversionResponse{}, fmt.Errorf("generate conversion ID: %w", err)
	}

	if err := os.MkdirAll(s.cfg.OutputDir, 0750); err != nil {
		return model.ConversionResponse{}, fmt.Errorf("create output directory: %w", err)
	}

	outputTemplate := filepath.Join(s.cfg.OutputDir, id+".%(ext)s")

	args := []string{"--no-warnings", "--no-playlist", "--ffmpeg-location",  s.cfg.FFmpegPath, "--output", outputTemplate}

		switch request.Format {
	case "mp3":
		bitrate := strings.TrimSuffix(request.Quality, "kbps") + "K"

		args = append(args,
			"--extract-audio",
			"--audio-format", "mp3",
			"--audio-quality", bitrate,
		)

	case "mp4":
		height := strings.TrimSuffix(request.Quality, "p")

		formatSelector := fmt.Sprintf(
			"bestvideo[height<=%s][ext=mp4]+bestaudio[ext=m4a]/"+
				"best[height<=%s][ext=mp4]/best[height<=%s]",
			height, height, height,
		)

		args = append(args,
			"--format", formatSelector,
			"--merge-output-format", "mp4",
		)
	}

	args = append(args, "--", request.URL)

	conversionCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)

	defer cancel()

	cmd := exec.CommandContext(conversionCtx, s.cfg.YtDlpPath, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return model.ConversionResponse{}, fmt.Errorf("conversion failed: %w", err)
	}

	return model.ConversionResponse{
		ID:     id,
		Status: model.ConversionStatusCompleted,
	}, nil
}

func generateID()(string, error){

}

func validateConversionRequest(request model.ConversionRequest) error {
	parsedURL, err := url.ParseRequestURI(request.URL)

	// valid url scheme
	if err != nil || parsedURL.Scheme != "https" || parsedURL.User != nil {
		return fmt.Errorf("%w: invalid URL", ErrInvalidConversionRequest)
	}

	// valid host/domain
	switch strings.ToLower(parsedURL.Hostname()) {
	case "youtube.com", "www.youtube.com", "m.youtube.com", "youtu.be":
		// Allowed hostname.
	default:
		return fmt.Errorf("%w: unsupported URL host", ErrInvalidConversionRequest)
	}

	switch request.Format {
	case "mp3":
		switch request.Quality {
		case "128kbps", "192kbps", "320kbps":
			return nil
		default:
			return fmt.Errorf("%w: unsupported MP3 quality", ErrInvalidConversionRequest)
		}

	case "mp4":
		switch request.Quality {
		case "360p", "480p", "720p", "1080p", "1440p", "2160p":
			return nil
		default:
			return fmt.Errorf("%w: unsupported MP4 quality", ErrInvalidConversionRequest)
		}

	default:
		return fmt.Errorf("%w: unsupported output format", ErrInvalidConversionRequest)
	}

}