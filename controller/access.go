package controller

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/utils"
)

type accessController struct {
	RemoteAddr string      `json:"remote_addr"`
	Expiry     int64       `json:"expiry"`
	Role       string      `json:"role"`
	Payload    interface{} `json:"payload"`
}

func (ctl *accessController) newToken(secret string) (string, error) {
	v, err := json.Marshal(ctl)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create token: build body failed: %s",
			err.Error(),
		)
	}

	var body map[string]interface{}
	if err = json.Unmarshal(v, &body); err != nil {
		return "", fmt.Errorf(
			"failed to create token: build body failed: %s",
			err.Error(),
		)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims(body)

	return token.SignedString([]byte(secret))
}

func (ctl *accessController) initByToken(token, secret string) error {
	t, err := jwt.Parse(token, func(t1 *jwt.Token) (interface{}, error) {
		if v, ok := t1.Method.(*jwt.SigningMethodHMAC); !ok || v == nil {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if !t.Valid {
		return fmt.Errorf("not a valid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("not valid claims")
	}

	d, err := json.Marshal(claims)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, ctl)
}

func verifyCSRFToken(tokenbyte, csrftoken []byte) (ok bool) {
	token := string(tokenbyte)

	return string(csrftoken) == token
}

func (ctl *accessController) refreshToken(expiry int64, secret string) (string, error) {
	ctl.Expiry = utils.Expiry(expiry)
	return ctl.newToken(secret)
}

func (ctl *accessController) genCSRFToken(token string) (ct string) {
	return token
}

func (ctl *accessController) verify(roles []string) error {
	if ctl.Expiry < utils.Now() {
		return fmt.Errorf("token is expired")
	}

	if !sets.NewString(roles...).Has(ctl.Role) {
		return fmt.Errorf("not allowed permissions")
	}

	return nil
}
