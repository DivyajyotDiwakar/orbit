package signin

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func SignInService(c context.Context, req SignInRequest) (string, *models.User, error) {
	user, err := repositories.GetUserFromEmail(c, req.EmailId)
	if err != nil {
		return "", nil, err
	}

	if hash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); hash != nil {
		return "", nil, errors.New("EmailId or Password is wrong..!!")
	}

	token, err := utils.CreateJwtToken(*user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
