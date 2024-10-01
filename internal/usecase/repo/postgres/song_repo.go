package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"songs/internal/models"
	"time"
)

type SongRepo struct {
	*sql.DB
}

type Songs interface {
	GetSongs(ctx context.Context, groupName string, releaseDate time.Time, offset, limit int) (map[string][]*models.Song, int, int, error)
	Add(ctx context.Context, groupName, title, text, link string, releaseDate time.Time) (int, error)
	Delete(ctx context.Context, ID int) (int, error)
	Update(ctx context.Context, ID int, title, text, link string, releaseDate time.Time) (int, error)
	GetLyric(ctx context.Context, ID int) (string, error)
}

func NewSongRepo(db *sql.DB) *SongRepo {
	return &SongRepo{db}
}

var _ Songs = (*SongRepo)(nil)

func (s SongRepo) GetSongs(ctx context.Context, groupName string, releaseDate time.Time, offset, limit int) (map[string][]*models.Song, int, int, error) {
	const op = "SongRepo - GetSongs"

	var params []interface{}

	query := "SELECT name, songs.title, songs.release_date, songs.link, songs.lyric FROM groups " +
		"LEFT JOIN songs ON groups.id = songs.group_id " +
		"WHERE 1=1"

	paramIdx := 1

	if groupName != "" {
		query += fmt.Sprintf(" AND name = $%d ", len(params)+1)
		//query += " AND name = $" + fmt.Sprint(paramIdx) + " "
		params = append(params, groupName)
		paramIdx++
	}

	var zeroTime time.Time

	if releaseDate != zeroTime {
		query += fmt.Sprintf(" AND songs.release_date = $%d ", len(params)+1)
		//query += " AND songs.release_date = $" + fmt.Sprint(paramIdx) + " "
		params = append(params, releaseDate)
		paramIdx++
	}

	query += " LIMIT $" + fmt.Sprint(paramIdx) + " OFFSET $" + fmt.Sprint(paramIdx+1)
	params = append(params, limit, offset)

	rows, err := s.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%s - s.QueryContext :%w", op, err)
	}

	defer rows.Close()

	songs := map[string][]*models.Song{}

	for rows.Next() {
		var song models.Song

		var name string

		rows.Scan(&name, &song.Title, &song.ReleaseDate, &song.Link, &song.Lyric)

		songs[name] = append(songs[name], &song)
	}

	return songs, offset, limit, nil
}

func (s SongRepo) Add(ctx context.Context, groupName, title, lyric, link string, releaseDate time.Time) (int, error) {
	const op = "SongRepo - Add"

	var groupID int

	query := "INSERT INTO groups(name)" +
		"VALUES($1) ON CONFLICT (name) DO UPDATE " +
		"SET name = groups.name " +
		"RETURNING id"

	err := s.QueryRowContext(ctx, query, groupName).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("%s - s.QueryRowContext - InsertGroup: %w", op, err)
	}

	query = "INSERT INTO songs(title, group_id, release_date, link, lyric)" +
		"VALUES($1, $2, $3, $4, $5) RETURNING ID"

	var songID int

	err = s.QueryRowContext(ctx, query, title, groupID, releaseDate, link, lyric).Scan(&songID)
	if err != nil {
		return 0, fmt.Errorf("%s - s.QueryRowContext - InsertSong: %w", op, err)
	}

	return songID, nil
}

func (s SongRepo) Delete(ctx context.Context, ID int) (int, error) {
	const op = "SongRepo - Delete"

	query := "DELETE FROM songs WHERE ID = $1"

	_, err := s.ExecContext(ctx, query, ID)
	if err != nil {
		return 0, fmt.Errorf("%s - s.ExecContext: %w", op, err)
	}

	return ID, nil
}

func (s SongRepo) Update(ctx context.Context, ID int, title, lyric, link string, releaseDate time.Time) (int, error) {
	const op = "SongRepo - Update"

	query := "UPDATE songs SET " +
		"title = COALESCE(NULLIF($1, ''), title), " +
		"release_date = COALESCE(NULLIF($2, $3)::timestamp, release_date), " +
		"link = COALESCE(NULLIF($4, ''), link), " +
		"lyric = COALESCE(NULLIF($5, ''), lyric) " +
		"WHERE id = $6"

	var zeroTime time.Time

	res, err := s.ExecContext(ctx, query, title, releaseDate, zeroTime, link, lyric, ID)
	if err != nil {
		return 0, fmt.Errorf("%s - ExecContext: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s - res.RowAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("%s - no song updated", op)
	}

	return ID, nil
}

func (s SongRepo) GetLyric(ctx context.Context, ID int) (string, error) {
	const op = "SongRepo - GeyLyric"

	query := "SELECT lyric FROM songs " +
		"WHERE id = $1 "

	rows, err := s.QueryContext(ctx, query, ID)
	if err != nil {
		return "", fmt.Errorf("%s - s.QueryContext: %w", op, err)
	}

	defer rows.Close()

	var verses string

	rows.Next()
	rows.Scan(&verses)

	return verses, err
}
