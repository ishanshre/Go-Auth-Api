package validators

import (
	"github.com/ishanshre/go-auth-api/internals/pkg/utils"
)

func ValidateUsername(username string) error {
	if err := utils.CheckUsername(username); err != nil {
		return err
	}
	return nil
}

func ValidatePassword(password string) error {
	if err := utils.CheckPassword(password); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	return utils.CheckEmail(email)
}
