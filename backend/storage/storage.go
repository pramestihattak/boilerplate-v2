package storage

import "context"

type AuthStorage interface {
	AuthWriter
	AuthReader
}

type AuthReader interface {
	CheckedUserForVerification(ctx context.Context, email, verificationToken string) (int, error)
	UserExist(ctx context.Context, email string) (int, error)
	Login(ctx context.Context, input *LoginInput) (*LoginOutput, error)
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

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	UserID   string
	Email    string
	FullName string
	Password string
}
