package db

import (
	"database/sql"
	"fmt"

	"github.com/OskolockKoli/url_shortener/internal/models"
)

type PostgreSQL struct {
	DB *sql.DB
}

func (pg *PostgreSQL) Save(link models.Link) error {
	query := `
    INSERT INTO links (short_id, url)
    VALUES ($1, $2)
    ON CONFLICT DO NOTHING
    RETURNING id
    `

	var id int64
	err := pg.DB.QueryRow(query, link.ShortID, link.URL).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to insert link: %w", err)
	}

	return nil
}

func (pg *PostgreSQL) GetByShortID(shortID string) (models.Link, error) {
	query := `
    SELECT short_id, url
    FROM links
    WHERE short_id = $1
    LIMIT 1
    `

	row := pg.DB.QueryRow(query, shortID)
	var link models.Link
	err := row.Scan(&link.ShortID, &link.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, fmt.Errorf("failed to get link by short ID: %w", err)
	}

	return link, nil
}
