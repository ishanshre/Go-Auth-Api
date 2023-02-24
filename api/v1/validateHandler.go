package v1

import (
	"github.com/ishanshre/go-auth-api/api/v1/models"
	"github.com/ishanshre/go-auth-api/internals/pkg/validators"
)

func validateCreateUser(user *models.RegisterUser) error {
	// if err := validators.ValidateUsername(user.Username); err != nil {
	// 	return err
	// }
	// if err := validators.ValidatePassword(user.Password); err != nil {
	// 	return err
	// }
	if err := validators.ValidateEmail(user.Email); err != nil {
		return err
	}
	return nil
}
