// Package apple handles the validation and authentication with
// Apple servers for Apple Sign In token validations.
package apple

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	validationEndpoint = "https://appleid.apple.com/auth/token"
	appleAudience      = "https://appleid.apple.com"
)

// AppleAuth is the contract for communication and validation of
// Apple user tokens.
type AppleAuth interface {
	// ValidateCode validates an authorization code returning the associated
	// Apple user and the generated tokens(access and refresh).
	ValidateCode(code string) (*TokenResponse, error)
}

type AppleUser struct{}

// TokenResponse response when validation was successfull.
type TokenResponse struct {
	// AccessToken (Reserved for future use) A token used to access allowed data.
	// Currently, no data set has been defined for access.
	AccessToken string
	// ExpiresIn the amount of time, in seconds, before the access token expires.
	ExpiresIn int
	// IDToken a JSON Web Token that contains the user’s identity information.
	IDToken string
	// RefreshToken The refresh token used to regenerate new access tokens.
	// Store this token securely on your server.
	RefreshToken string
	// TokenType the type of access token.
	TokenType string
}

type appleAuth struct {
	AppID      string
	TeamID     string
	KeyID      string
	KeyContent []byte
}

// Setup and return a new AppleAuth for validation of tokens.
func New(appID, teamID, keyID, keyPath string) (*appleAuth, error) {
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &appleAuth{
		KeyID:      keyID,
		TeamID:     teamID,
		AppID:      appID,
		KeyContent: keyContent,
	}, nil
}

func (a *appleAuth) clientSecret() (string, error) {
	block, _ := pem.Decode(a.KeyContent)
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Second * 15776999).Unix(),
		Issuer:    a.TeamID,
		Subject:   a.AppID,
		Audience:  appleAudience,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = a.KeyID
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return clientSecret, nil
}

func (a *appleAuth) ValidateCode(code string) (*TokenResponse, error) {
	clientSecret, err := a.clientSecret()
	if err != nil {
		return nil, err
	}
	var formQuery url.Values
	formQuery.Add("client_id", a.AppID)
	formQuery.Add("client_secret", clientSecret)
	formQuery.Add("code", code)
	formQuery.Add("grant_type", "authorization_code")
	res, err := http.PostForm(validationEndpoint, formQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		var errorResponseBody struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(res.Body).Decode(&errorResponseBody); err != nil {
			return nil, err
		}
		switch errorResponseBody.Error {
		case string(ErrorResponseTypeInvalidScope):
			return nil, ErrorResponseInvalidScope
		case string(ErrorResponseTypeUnsupportedGrantType):
			return nil, ErrorResponseUnsupportedGrantType
		case string(ErrorResponseTypeUnauthorizedClient):
			return nil, ErrorResponseUnauthorizedClient
		case string(ErrorResponseTypeInvalidGrant):
			return nil, ErrorResponseInvalidGrant
		case string(ErrorResponseTypeInvalidClient):
			return nil, ErrorResponseInvalidClient
		case string(ErrorResponseTypeInvalidRequest):
			return nil, ErrorResponseInvalidRequest
		default:
			return nil, fmt.Errorf("unrecognized response error: %s", errorResponseBody.Error)
		}
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}
	return &tokenResponse, nil
}
