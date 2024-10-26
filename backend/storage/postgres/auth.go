package postgres

import (
	"context"
	"database/sql"

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

	userExistSQL = `
		SELECT 
			COUNT(user_id)
		FROM users
		WHERE email = $1 LIMIT 1;
	`

	checkedUserVerificationExistSQL = `
		SELECT 
			COUNT(user_id)
		FROM users
		WHERE email = $1 and verification_token = $2 LIMIT 1;
	`

	verifyUserSQL = `
		UPDATE users
			SET verified = true, verification_token = ''
		WHERE email = $1
		RETURNING user_id
	`

	loginSQL = `
		SELECT 
			user_id,
			full_name,
			email,
			password,
			verified
		FROM users
		WHERE email = $1;
	`
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

func (s *Storage) UserExist(ctx context.Context, email string) (int, error) {
	var exist int
	err := s.db.GetContext(ctx, &exist, userExistSQL, email)
	if err != nil {
		return 0, err
	}

	return exist, nil
}

func (s *Storage) CheckedUserForVerification(ctx context.Context, email, verificationToken string) (int, error) {
	var exist int
	err := s.db.GetContext(ctx, &exist, checkedUserVerificationExistSQL, email, verificationToken)
	if err != nil {
		return 0, err
	}

	return exist, nil
}

func (s *Storage) VerifyUser(ctx context.Context, email string) (string, error) {
	txn, err := s.db.Begin()
	if err != nil {
		return "", errors.Wrap(err, "failed to verify, can't init transaction")
	}

	var id string
	if err := txn.QueryRowContext(ctx, verifyUserSQL,
		email,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		if err := txn.Rollback(); err != nil {
			return "", errors.Wrap(err, "failed to verify, can't rollback transaction")
		}
		return "", errors.Wrap(err, "failed to verify")
	}

	if err := txn.Commit(); err != nil {
		return "", errors.Wrap(err, "failed to verify, can't commit transaction")
	}
	return id, nil
}

func (s *Storage) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	var output LoginOutput
	if err := s.db.QueryRowContext(ctx, loginSQL, input.Email).Scan(
		&output.UserID,
		&output.FullName,
		&output.Email,
		&output.Password,
		&output.Verified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &output, nil
}
