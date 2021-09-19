package backend

import (
	"fmt"
	"net/url"
	"strings"
)

type AccountParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register creates a new user and returns their session token.
func Register(email, password string) (string, error) {
	body := AccountParams{
		Email:    email,
		Password: password,
	}
	var token string
	if err := Post("", "/register", body, &token); err != nil {
		return "", err
	}
	return token, nil
}

// Login authenticate a user and returns their session token
func Login(email, password string) (string, error) {
	body := AccountParams{
		Email:    email,
		Password: password,
	}
	var token string
	if err := Post("", "/login", body, &token); err != nil {
		return "", err
	}

	return token, nil
}

func SetPassword(token, email, oldPassword, newPassword string) error {
	var body = new(struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	})

	body.Email = email
	body.OldPassword = oldPassword
	body.NewPassword = newPassword

	var status bool
	if err := Post(token, "/user/changepw", body, &status); err != nil {
		return err
	}
	return nil
}

func GetPasswordResetCode(token, email string) (string, error) {
	qs := url.Values{}
	qs.Add("e", email)

	var code string
	path := fmt.Sprintf("/password/resetcode?%s", qs.Encode())
	if err := Get(token, path, &code); err != nil {
		return "", err
	}
	return code, nil
}

func ResetPassword(email, code, password string) error {
	var body = new(struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	})
	body.Email = strings.ToLower(email)
	body.Code = code
	body.Password = password

	var status bool
	if err := Post("", "/password/reset", body, &status); err != nil {
		return err
	} else if !status {
		return fmt.Errorf("unable to reset password")
	}
	return nil
}

// SudoGetToken returns a token from an AccountID
// This is useful when performing creation that documents needs
// to be attached to a specific account id and therefor the SudoCreate
// would not work on those cases
func SudoGetToken(token, accountID string) (string, error) {
	var tok string
	if err := Get(token, "/sudogettoken/"+accountID, &tok); err != nil {
		return "", err
	}
	return tok, nil
}
