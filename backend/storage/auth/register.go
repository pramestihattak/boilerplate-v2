package auth_storage

import (
	"context"

	"github.com/pkg/errors"
)

// sql queries
var (
	registerSQL = `
		INSERT INTO users (
			full_name,
			email,
			password,
			verification_token
		) VALUES (
				$1, $2, $3, $4
		) RETURNING user_id`
)

func (s *Storage) Register(ctx context.Context, reg Register) (string, error) {
	txn, err := s.db.Begin()
	if err != nil {
		return "", errors.Wrap(err, "failed to register, can't init transaction")
	}

	var id string
	if err := txn.QueryRowContext(ctx, registerSQL,
		reg.FullName,
		reg.Email,
		reg.Password,
		reg.VerificationToken,
	).Scan(&id); err != nil {
		if err := txn.Rollback(); err != nil {
			return "", errors.Wrap(err, "failed to register, can't rollback transaction")
		}
		return "", errors.Wrap(err, "failed to register")
	}

	if err := txn.Commit(); err != nil {
		return "", errors.Wrap(err, "failed to register, can't commit transaction")
	}
	return id, nil
}
