package backend

import (
	"fmt"
	"net/url"
)

// NewSystemAccount initiates the StaticBackend account creation process.
func NewSystemAccount(email string) (string, error) {
	var stripeURL string

	path := fmt.Sprintf("/account/init?email=%s", url.QueryEscape(email))
	if err := Get("", path, &stripeURL); err != nil {
		return "", err
	}
	return stripeURL, nil
}

// NewSystemAccountData when bypassing Stripe, this struct will be returned
// when creating a new system account.
type NewSystemAccountData struct {
	PublicKey     string `json:"pk"`
	RootToken     string `json:"rootToken"`
	AdminPassword string `json:"pw"`
}

func NewSystemAccountBypassStripe(email, bypassFlag string) (data NewSystemAccountData, err error) {
	q := url.Values{}
	q.Add("email", email)
	q.Add("ui", "true")
	q.Add("x", bypassFlag)

	u := fmt.Sprintf("/account/init?%s", q.Encode())
	fmt.Println("DEBUG", u)
	err = Get("", u, &data)
	return
}
