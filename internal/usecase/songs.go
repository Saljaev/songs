package usecase

import (
	"context"
	"fmt"
	"songs/internal/models"
	"songs/internal/usecase/repo/postgres"
	"strings"
	"time"
)

type SongUseCase struct {
	song postgres.SongRepo
}

func NewSongUseCase(repo postgres.SongRepo) *SongUseCase {
	return &SongUseCase{song: repo}
}

func (s *SongUseCase) Add(ctx context.Context, groupName, title, releaseDate, lyric, link string) (int, error) {
	const op = "SongUseCase - Add"

	release, err := time.Parse(time.DateOnly, releaseDate)
	if err != nil {
		return 0, fmt.Errorf("%s - failed to decode releaseDate: %w", op, err)
	}

	ID, err := s.song.Add(ctx, groupName, title, lyric, link, release)
	if err != nil {
		return 0, fmt.Errorf("%s - failed to add song: %w", op, err)
	}

	return ID, nil
}

func (s *SongUseCase) GetSongs(ctx context.Context, groupName, releaseDate string, offset, limit int) (map[string][]*models.Song, int, int, error) {
	const op = "SongUseCase - GetSongs"

	var release time.Time

	if releaseDate != "" {
		var err error

		release, err = time.Parse(time.DateOnly, releaseDate)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("%s - failed to decode releaseDate: %w", op, err)
		}
	}

	songs, offset, limit, err := s.song.GetSongs(ctx, groupName, release, offset, limit)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%s - failed to get songs: %w", op, err)
	}

	return songs, offset, limit, nil
}

func (s *SongUseCase) Delete(ctx context.Context, ID int) (int, error) {
	const op = "SongUseCase - Delete"

	ID, err := s.song.Delete(ctx, ID)
	if err != nil {
		return 0, fmt.Errorf("%s - failed to delete song by id: %w", op, err)
	}

	return ID, nil
}

func (s *SongUseCase) Update(ctx context.Context, ID int, title, lyric, releaseDate, link string) (int, error) {
	const op = "SongUseCase - Update"

	var release time.Time

	if releaseDate != "" {
		var err error

		release, err = time.Parse(time.DateOnly, releaseDate)
		if err != nil {
			return 0, fmt.Errorf("%s - failed to decode releaseDate: %w", op, err)
		}
	}

	ID, err := s.song.Update(ctx, ID, title, lyric, link, release)
	if err != nil {
		return 0, fmt.Errorf("%s - failed to update song: %w", op, err)
	}

	return ID, err

}

func (s *SongUseCase) GetLyric(ctx context.Context, ID, offset, limit int) ([]string, error) {
	const op = "SongUseCase - GetLyric"

	lyric, err := s.song.GetLyric(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("%s - failed to get lyric: %w", op, err)
	}

	verses := strings.Split(lyric, "\n\n")

	if offset > len(verses) {
		offset = len(verses)
	}

	if limit+offset > len(verses) {
		limit = len(verses)
	}

	return verses[offset : offset+limit], nil
}
