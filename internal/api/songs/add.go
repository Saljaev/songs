package songs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"songs/internal/api"
	"unicode/utf8"
)

type AddReq struct {
	Group string `json:"group"`
	Song  string `json:"song"`
	SongDetail
}

type SongDetail struct {
	ReleaseDate string `json:"release_date"`
	Lyric       string `json:"lyric"`
	Link        string `json:"link"`
}

func (req AddReq) Validate(_ *api.Context) error {
	if utf8.RuneCountInString(req.Group) == 0 {
		return errors.New("zero group name")
	}
	if utf8.RuneCountInString(req.Song) == 0 {
		return errors.New("zero song name")
	}

	return nil
}

type SongIDResp struct {
	ID int `json:"id"`
}

func (s *SongsHandler) Add(_ *api.Context, req *AddReq) (*SongIDResp, int) {
	s.log.Debug("received request", slog.Any("request", req))

	songDetail, err := AddExternal(req.Group, req.Song)
	if err != nil {
		s.log.Error("failed to add get song detail", slog.String("error", err.Error()))

		return &SongIDResp{
			ID: 0,
		}, http.StatusInternalServerError
	}

	s.log.Debug("fetched song details from external service", slog.Any("songDetail", songDetail))

	ID, err := s.songs.Add(context.Background(), req.Group, req.Song, songDetail.ReleaseDate, songDetail.Link, songDetail.Lyric)
	if err != nil {
		s.log.Error(err.Error())

		return &SongIDResp{
			ID: 0,
		}, http.StatusInternalServerError
	}

	resp := &SongIDResp{ID: ID}

	s.log.Info("successfully added song", slog.Int("ID", ID), slog.String("group", req.Group), slog.String("song", req.Song))

	return resp, http.StatusOK
}

func AddExternal(group, song string) (*SongDetail, error) {
	/*
		Без имитации код:
		url := fmt.Sprintf("http://example/info?group=%s&song=%s", group, song)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("http.Get: %w", err)
		}

		defer resp.Body.Close()

		var data *SongDetail

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode song detail: %w", err)
		}

		return data, nil
	*/

	//Имитация выполнения
	return &SongDetail{
		ReleaseDate: "2006-10-10",
		Lyric:       "hello \n\n world \n\n a'am",
		Link:        "https://youtube.com",
	}, nil
}
