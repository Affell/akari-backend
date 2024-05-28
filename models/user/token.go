package user

import (
	"errors"
	"time"

	"github.com/fatih/structs"
	"github.com/google/uuid"
)

var tokens map[string]UserToken = make(map[string]UserToken)

type UserToken struct {
	TokenID   string    `json:"token_id" structs:"-"`
	ID        int64     `json:"id" structs:"id"`
	Username  string    `json:"username" structs:"username"`
	Email     string    `json:"email" structs:"email"`
	CreatedAt time.Time `json:"created_at" structs:"-"`
}

func (token UserToken) IsNil() bool {
	return token.TokenID == ""
}

func (token UserToken) ToUserData() map[string]interface{} {
	return structs.Map(token)
}

func (userToken *UserToken) Store() (tokenID string) {

	if userToken.TokenID == "" {
		tokenID = uuid.New().String()
		userToken.TokenID = tokenID
	} else {
		tokenID = userToken.TokenID
	}

	tokens[userToken.TokenID] = UserToken{
		ID:        userToken.ID,
		TokenID:   userToken.TokenID,
		Username:  userToken.Username,
		Email:     userToken.Email,
		CreatedAt: userToken.CreatedAt,
	}

	return
}

func (userToken *UserToken) Update(u User) {
	userToken.ID = u.ID
	userToken.Username = u.Username
	userToken.Email = u.Email
}

func GetUserToken(tokenID string) (userToken UserToken, err error) {

	userToken, ok := tokens[tokenID]
	if !ok {
		err = errors.New("incorrect token")
	}

	return
}

func RevokeUserToken(tokenID string) {
	delete(tokens, tokenID)
}
