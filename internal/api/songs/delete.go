package songs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"songs/internal/api"
)

type DeleteReq struct {
	ID int `query:"id"`
}

func (req DeleteReq) Validate(_ *api.Context) error {
	if req.ID <= 0 {
		return errors.New("invalid input id")
	}

	return nil
}

func (s *SongsHandler) Delete(_ *api.Context, req *DeleteReq) (*SongIDResp, int) {
	s.log.Debug("received request", slog.Any("request", req))

	ID, err := s.songs.Delete(context.Background(), req.ID)
	if err != nil {
		s.log.Error("failed to delete song", slog.Any("ID", req.ID), slog.String("error", err.Error()))

		return &SongIDResp{
			ID: 0,
		}, http.StatusInternalServerError
	}

	resp := &SongIDResp{ID: ID}

	s.log.Info("successfully deleted song", slog.Int("song ID", ID))

	return resp, http.StatusOK
}
