package storage

import "context"

type AuthStorage interface {
	AuthWriter
	AuthReader
}

type AuthReader interface {
	CheckedUserForVerification(ctx context.Context, email, verificationToken string) (int, error)
	UserExist(ctx context.Context, email string) (int, error)
}

type AuthWriter interface {
	Register(ctx context.Context, reg Register) (string, error)
	VerifyUser(ctx context.Context, email string) (string, error)
}

type Register struct {
	UserID            string `db:"user_id"`
	FullName          string `db:"full_name"`
	Email             string `db:"email"`
	Password          string `db:"password"`
	VerificationToken string `db:"verification_token"`
}
