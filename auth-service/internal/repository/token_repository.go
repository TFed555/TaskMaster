package repository

import (
	"auth-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *models.RefreshToken) error {

	stmt, err := r.db.PrepareNamed(`
    INSERT INTO refresh_tokens (user_id, token, expires_at)
    VALUES (:user_id, :token, :expires_at)
    ON CONFLICT (user_id) DO UPDATE SET
        token = EXCLUDED.token,
        expires_at = EXCLUDED.expires_at
    RETURNING id, expires_at
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.Get(token, token)
}

func (r *TokenRepository) GetToken(token string) (*models.RefreshToken, error) {
	query := `SELECT * FROM refresh_tokens
        WHERE token = $1
        LIMIT 1`
	var result models.RefreshToken
	err := r.db.Get(&result, query, token)
	return &result, err
}

func (r *TokenRepository) DeleteToken(refreshToken string) (bool, error) {
	query := `DELETE FROM refresh_tokens WHERE token = $1`

	res, err := r.db.Exec(query, refreshToken)

	rows, err := res.RowsAffected()
    if err != nil {
        return false, err
    }
    
    return rows > 0, nil
}