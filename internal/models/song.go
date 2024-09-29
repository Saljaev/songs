package models

import (
	"strings"
	"time"
)

type (
	Song struct {
		Title string
		SongDetail
	}

	SongDetail struct {
		ReleaseDate time.Time
		Link        string
		Lyric       string
	}
)

func (s *Song) SplitLyric() []string {
	verses := strings.Split(s.Lyric, "\n\n")

	return verses
}
