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
