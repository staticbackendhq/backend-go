package backend

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
