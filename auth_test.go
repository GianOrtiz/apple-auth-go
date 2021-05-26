package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Implementation of a mocked HTTP client.
type MockedHTTPClient struct {
	mock.Mock
}

// Mocked function PostForm that does not call any server, just return the expected response.
func (m *MockedHTTPClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := m.Mock.Called(url, data)

	resArg := args.Get(0)
	resp, ok := resArg.(*http.Response)
	if !ok {
		return nil, errors.New("first parameter should be of type *http.Response")
	}

	err = args.Error(1)
	return resp, err
}

const mockClientSecret = "client-secret"

func TestValidateRequest(t *testing.T) {
	form := make(url.Values)

	tokenResponse := TokenResponse{}
	tokenResponseBody, _ := json.Marshal(tokenResponse)
	mockedHTTPClient := new(MockedHTTPClient)
	mockedHTTPClient.On("PostForm", validationEndpoint, form).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tokenResponseBody)),
		},
		nil,
	)

	auth := appleAuth{
		AppID:      "appID",
		TeamID:     "teamID",
		KeyID:      "keyID",
		KeyContent: []byte{},
		httpClient: mockedHTTPClient,
	}
	res, err := auth.validateRequest(form)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func TestValidateCode(t *testing.T) {
	code := "apple-authorization-code"

	tokenResponse := TokenResponse{}
	tokenResponseBody, _ := json.Marshal(tokenResponse)
	mockedHTTPClient := new(MockedHTTPClient)

	auth := appleAuth{
		AppID:      "appID",
		TeamID:     "teamID",
		KeyID:      "keyID",
		KeyContent: []byte{},
		httpClient: mockedHTTPClient,
	}
	reqForm := make(url.Values)
	reqForm.Add("client_id", auth.AppID)
	reqForm.Add("client_secret", mockClientSecret)
	reqForm.Add("code", code)
	reqForm.Add("grant_type", "authorization_code")
	mockedHTTPClient.On("PostForm", validationEndpoint, reqForm).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tokenResponseBody)),
		},
		nil,
	)

	res, err := auth.validateCode(mockClientSecret, code)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func TestValidateCodeWithRedirectURI(t *testing.T) {
	code := "apple-authorization-code"
	redirectURI := "https://saladeestar.app/apple"

	tokenResponse := TokenResponse{}
	tokenResponseBody, _ := json.Marshal(tokenResponse)
	mockedHTTPClient := new(MockedHTTPClient)

	auth := appleAuth{
		AppID:      "appID",
		TeamID:     "teamID",
		KeyID:      "keyID",
		KeyContent: []byte{},
		httpClient: mockedHTTPClient,
	}
	reqForm := make(url.Values)
	reqForm.Add("client_id", auth.AppID)
	reqForm.Add("client_secret", mockClientSecret)
	reqForm.Add("code", code)
	reqForm.Add("grant_type", "authorization_code")
	reqForm.Add("redirect_uri", redirectURI)
	mockedHTTPClient.On("PostForm", validationEndpoint, reqForm).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tokenResponseBody)),
		},
		nil,
	)

	res, err := auth.validateCodeWithRedirectURI(mockClientSecret, code, redirectURI)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func TestValidateRefreshToken(t *testing.T) {
	refreshToken := "refresh-token-as-jwt"

	tokenResponse := TokenResponse{}
	tokenResponseBody, _ := json.Marshal(tokenResponse)
	mockedHTTPClient := new(MockedHTTPClient)

	auth := appleAuth{
		AppID:      "appID",
		TeamID:     "teamID",
		KeyID:      "keyID",
		KeyContent: []byte{},
		httpClient: mockedHTTPClient,
	}
	reqForm := make(url.Values)
	reqForm.Add("client_id", auth.AppID)
	reqForm.Add("client_secret", mockClientSecret)
	reqForm.Add("refresh_token", refreshToken)
	reqForm.Add("grant_type", "refresh_token")
	mockedHTTPClient.On("PostForm", validationEndpoint, reqForm).Return(
		&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tokenResponseBody)),
		},
		nil,
	)

	res, err := auth.validateRefreshToken(mockClientSecret, refreshToken)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}
