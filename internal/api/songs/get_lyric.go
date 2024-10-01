package songs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"songs/internal/api"
)

type LyricReq struct {
	ID     int `query:"id"`
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

func (req LyricReq) Validate(_ *api.Context) error {
	if req.ID <= 0 {
		return errors.New("invalid input id")
	}
	if req.Offset < 0 {
		return errors.New("invalid offset value")
	}
	if req.Limit < 1 {
		return errors.New("invalid limit value")
	}
	return nil
}

type LyricResp struct {
	Verses []string `json:"verses"`
}

func (s *SongsHandler) GetLyric(_ *api.Context, req *LyricReq) (*LyricResp, int) {
	s.log.Debug("received request", slog.Any("request", req))

	verses, err := s.songs.GetLyric(context.Background(), req.ID, req.Offset, req.Limit)
	if err != nil {
		s.log.Error("failed to get song's text ", slog.Any("ID", req.ID), slog.String("error", err.Error()))

		return &LyricResp{
			Verses: []string{},
		}, http.StatusInternalServerError
	}

	resp := &LyricResp{Verses: verses}

	s.log.Info("successfully get song's text", slog.Any("song ID", req.ID))

	return resp, http.StatusOK
}
