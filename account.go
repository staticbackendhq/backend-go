package backend

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// AccountParams represents a new StaticBackend account
type AccountParams struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID string `json:"accountId,omitempty"`
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
	return LoginForAccount(email, password, "")
}

// LoginForAccount authenticate a user ensuring they're part of
// an account and returns their session token
func LoginForAccount(email, password, accountID string) (string, error) {
	body := AccountParams{
		Email:     email,
		Password:  password,
		AccountID: accountID,
	}
	var token string
	if err := Post("", "/login", body, &token); err != nil {
		return "", err
	}

	return token, nil
}

// EmailExists checks whether an email address is already registered.
func EmailExists(email string) (bool, error) {
	qs := url.Values{}
	qs.Add("e", email)

	var exists bool
	path := fmt.Sprintf("/email?%s", qs.Encode())
	if err := Get("", path, &exists); err != nil {
		return false, err
	}
	return exists, nil
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

// SetRole changes the role of a user in an account.
func SetRole(token, accountID, email string, role int) error {
	var body = new(struct {
		AccountID string `json:"accountId"`
		Email     string `json:"email"`
		Role      int    `json:"role"`
	})

	body.AccountID = accountID
	body.Email = email
	body.Role = role

	var status bool
	if err := Post(token, "/setrole", body, &status); err != nil {
		return err
	} else if !status {
		return fmt.Errorf("unable to set role")
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

// User represents a user record returned by root-token account endpoints.
type User struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	Role      int       `json:"role"`
	Created   time.Time `json:"created"`
}

// Me returns the current user matching this session token
// This is the only way to get the user's role, account/user ids and email.
func Me(token string) (me CurrentUser, err error) {
	err = Get(token, "/me", &me)
	return
}

// ChangeEmail changes the authenticated user's email address.
func ChangeEmail(token, email string) error {
	body := new(struct {
		Email string `json:"email"`
	})
	body.Email = email

	var status bool
	if err := Post(token, "/me/email", body, &status); err != nil {
		return err
	} else if !status {
		return fmt.Errorf("unable to change email")
	}
	return nil
}

// Users returns all users for the account linked with this token
func Users(token string) ([]CurrentUser, error) {
	var users []CurrentUser
	if err := Get(token, "/account/users", &users); err != nil {
		return nil, err
	}

	return users, nil
}

// AccountUser represents a cross-account membership for a user.
type AccountUser struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	AccountID string `json:"accountId"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	Token     string `json:"token"`
}

// UserAccountEntry represents one account (home or associated) returned by SudoGetUserAccounts.
type UserAccountEntry struct {
	AccountID string `json:"accountId"`
	Role      int    `json:"role"`
	Home      bool   `json:"home"`
	Token     string `json:"token,omitempty"`
}

// ListAssociations returns all cross-account memberships for the current user.
func ListAssociations(token string) ([]AccountUser, error) {
	var associations []AccountUser
	if err := Get(token, "/account/associations", &associations); err != nil {
		return nil, err
	}
	return associations, nil
}

// PromoteUser promotes the current user to have their own home account,
// preserving their existing membership as a cross-account association.
// Returns the new session token for the promoted account.
func PromoteUser(token string) (string, error) {
	var newToken string
	if err := Post(token, "/account/promote", nil, &newToken); err != nil {
		return "", err
	}
	return newToken, nil
}

// SudoGetUserAccounts returns all account IDs (home + associations) for a given email.
// Requires a root token.
func SudoGetUserAccounts(token, email string) ([]UserAccountEntry, error) {
	path := fmt.Sprintf("/account/user-accounts?email=%s", url.QueryEscape(email))
	var entries []UserAccountEntry
	if err := Get(token, path, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// SudoGetAuthTokenByUserID returns an authentication token for a specific user in an account.
// Requires a root token.
func SudoGetAuthTokenByUserID(token, accountID, userID string) (string, error) {
	path := fmt.Sprintf("/sudogetauthtokenbyuserid/%s/%s", accountID, userID)
	var authToken string
	if err := Get(token, path, &authToken); err != nil {
		return "", err
	}
	return authToken, nil
}

// SudoGetUserByID returns a specific user in an account.
// Requires a root token.
func SudoGetUserByID(token, accountID, userID string) (User, error) {
	path := fmt.Sprintf("/sudogetuserbyid/%s/%s", accountID, userID)
	var user User
	if err := Get(token, path, &user); err != nil {
		return User{}, err
	}
	return user, nil
}
