package backend

import (
	"fmt"
	"net/url"
	"strings"
)

// AccountParams represents a new StaticBackend account
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

// SetPassword changes the password of a user
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

// GetPasswordResetCode returns a unique code for a user to change their password
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

// ResetPassword changes user password using a unique code
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

// AddUser adds a user into the same account as token
func AddUser(token, email, password string) (CurrentUser, error) {
	body := AccountParams{
		Email:    email,
		Password: password,
	}
	var u CurrentUser
	err := Post(token, "/account/users", body, &u)
	return u, err
}

// RemoveUser removes a user from same account as token.
// Token must have a higher level of permission (role) than deleted user
func RemoveUser(token, userID string) error {
	uri := fmt.Sprintf("/account/users/%s", userID)
	return Del(token, uri)
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

// CurrentUser used to access current user's important information
type CurrentUser struct {
	AccountID string `json:"accountId"`
	UserID    string `json:"id"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
}

// Me returns the current user matching this session token
// This is the only way to get the user's role, account/user ids and email.
func Me(token string) (me CurrentUser, err error) {
	err = Get(token, "/me", &me)
	return
}

// Users returns all users for the account linked with this token
func Users(token string) ([]CurrentUser, error) {
	var users []CurrentUser
	if err := Get(token, "/account/users", &users); err != nil {
		return nil, err
	}

	return users, nil
}
