package otf

import (
	"fmt"
)

// Token is a user session
type Token struct {
	ID string `db:"token_id"`

	Token string

	Timestamps

	Description string

	// Token belongs to a user
	UserID string
}

func NewToken(uid, description string) (*Token, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	session := Token{
		ID:          NewID("ut"),
		Token:       token,
		Timestamps:  NewTimestamps(),
		Description: description,
		UserID:      uid,
	}

	return &session, nil
}