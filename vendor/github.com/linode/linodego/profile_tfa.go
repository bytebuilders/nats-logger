package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// TwoFactorSecret contains fields returned by CreateTwoFactorSecret
type TwoFactorSecret struct {
	Expiry *time.Time `json:"expiry"`
	Secret string     `json:"secret"`
}

// ConfirmTwoFactorOptions contains fields used by ConfirmTwoFactor
type ConfirmTwoFactorOptions struct {
	TFACode string `json:"tfa_code"`
}

// ConfirmTwoFactorResponse contains fields returned by ConfirmTwoFactor
type ConfirmTwoFactorResponse struct {
	Scratch string `json:"scratch"`
}

func (s *TwoFactorSecret) UnmarshalJSON(b []byte) error {
	type Mask TwoFactorSecret

	p := struct {
		*Mask
		Expiry *parseabletime.ParseableTime `json:"expiry"`
	}{
		Mask: (*Mask)(s),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	s.Expiry = (*time.Time)(p.Expiry)

	return nil
}

// CreateTwoFactorSecret generates a Two Factor secret for your User.
func (c *Client) CreateTwoFactorSecret(ctx context.Context) (*TwoFactorSecret, error) {
	e := "profile/tfa-enable"
	req := c.R(ctx).SetResult(&TwoFactorSecret{})
	r, err := coupleAPIErrors(req.Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*TwoFactorSecret), nil
}

// DisableTwoFactor disables Two Factor Authentication for your User.
func (c *Client) DisableTwoFactor(ctx context.Context) error {
	e := "profile/tfa-disable"
	_, err := coupleAPIErrors(c.R(ctx).Post(e))
	return err
}

// ConfirmTwoFactor confirms that you can successfully generate Two Factor codes and enables TFA on your Account.
func (c *Client) ConfirmTwoFactor(ctx context.Context, opts ConfirmTwoFactorOptions) (*ConfirmTwoFactorResponse, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	e := "profile/tfa-enable-confirm"
	req := c.R(ctx).SetResult(&ConfirmTwoFactorResponse{}).SetBody(string(body))
	r, err := coupleAPIErrors(req.Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*ConfirmTwoFactorResponse), nil
}
