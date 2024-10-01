package songs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"songs/internal/api"
	"time"
)

type SongReq struct {
	Group       string `json:"group,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Offset      int    `query:"offset"`
	Limit       int    `query:"limit"`
}

func (req SongReq) Validate(_ *api.Context) error {
	if req.Offset < 0 {
		return errors.New("invalid offset value")
	}
	if req.Limit < 1 {
		return errors.New("invalid limit value")
	}
	return nil
}

type SongResp struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Songs  map[string][]*SongData
}

type SongData struct {
	Title string `json:"title"`
	SongDetail
}

func (s *SongsHandler) GetSongs(_ *api.Context, req *SongReq) (*SongResp, int) {
	s.log.Debug("received request", slog.Any("request", req))

	songs, offset, limit, err := s.songs.GetSongs(context.Background(), req.Group, req.ReleaseDate, req.Offset, req.Limit)
	if err != nil {
		s.log.Error("failed to get songs", slog.String("error", err.Error()))

		return &SongResp{}, http.StatusInternalServerError
	}

	respSongs := map[string][]*SongData{}

	for k, v := range songs {
		for i := range v {
			song := SongData{
				Title: v[i].Title,

				SongDetail: SongDetail{
					v[i].ReleaseDate.Format(time.DateOnly),
					v[i].Lyric,
					v[i].Link,
				},
			}

			respSongs[k] = append(respSongs[k], &song)
		}
	}

	resp := &SongResp{
		Offset: offset,
		Limit:  limit,
		Songs:  respSongs,
	}

	resp.Offset = offset
	resp.Limit = limit

	s.log.Info("successfully get song", slog.Any("count", len(respSongs)))

	return resp, http.StatusOK
}
