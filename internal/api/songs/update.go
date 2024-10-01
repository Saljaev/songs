package songs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"songs/internal/api"
	"unicode/utf8"
)

type UpdateReq struct {
	ID          int    `query:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_Date"`
	Lyric       string `json:"lyric"`
	Link        string `json:"link"`
}

func (req UpdateReq) Validate(_ *api.Context) error {
	if req.ID <= 0 {
		return errors.New("invalid input id")
	}

	if utf8.RuneCountInString(req.Title) == 0 && utf8.RuneCountInString(req.ReleaseDate) == 0 &&
		utf8.RuneCountInString(req.Lyric) == 0 && utf8.RuneCountInString(req.Link) == 0 {
		return errors.New("zero value to update")
	}

	return nil
}

func (s *SongsHandler) Update(_ *api.Context, req *UpdateReq) (*SongIDResp, int) {
	s.log.Debug("received request", slog.Any("request", req))

	resp := &SongIDResp{}

	ID, err := s.songs.Update(context.Background(), req.ID, req.Title, req.Lyric, req.ReleaseDate, req.Link)
	if err != nil {
		s.log.Error("failed to update song", slog.Any("ID", req.ID), slog.String("error", err.Error()))

		return &SongIDResp{
			ID: 0,
		}, http.StatusInternalServerError
	}

	resp.ID = ID

	s.log.Info("successfully updated song", slog.Any("song ID", resp.ID))

	return resp, http.StatusOK
}
