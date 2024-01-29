package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Pastebin struct {
	Text        string
	OnlyOne     bool
	AliasForDel string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := storagePath
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePastebin(text string, alias string, aliasForDel string, onlyOne bool) error {
	const op = "internal.storage.postgres.SavePastebin"

	stmt, err := s.db.Prepare("INSERT INTO pastebin(text, alias, alias_for_del, only_one) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(text, alias, aliasForDel, onlyOne)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DelPastebin(aliasForDel string) error {
	const op = "internal.storage.postgres.DelPastebin"

	stmt, err := s.db.Prepare("UPDATE pastebin SET text = '[DEL]' || text, alias = alias || 'deletedByUser', alias_for_del = alias_for_del || 'deletedByUser', deleted_at = CURRENT_TIMESTAMP WHERE alias_for_del = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(aliasForDel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		// Обработка ситуации, когда данные не найдены
		return fmt.Errorf("%s: no data found for alias %s", op, aliasForDel)
	}

	return nil
}

func (s *Storage) ReadPastebin(alias string) (*Pastebin, error) {
	const op = "internal.storage.postgres.ReadPastebin"

	var p Pastebin

	stmt, err := s.db.Prepare("SELECT text, only_one, alias_for_del FROM pastebin WHERE alias= $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(alias).Scan(&p.Text, &p.OnlyOne, &p.AliasForDel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &p, nil
}
