package songs

import (
	"log/slog"
	"songs/internal/usecase"
)

type SongsHandler struct {
	songs usecase.SongUseCase
	log   *slog.Logger
}

func NewSongsHandler(s usecase.SongUseCase, l *slog.Logger) *SongsHandler {
	return &SongsHandler{
		songs: s,
		log:   l,
	}
}
