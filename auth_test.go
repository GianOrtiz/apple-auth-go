package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedHTTPClient struct {
	mock.Mock
}

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

var (
	mockedHTTPClient = new(MockedHTTPClient)
	auth             = appleAuth{
		AppID:      "appID",
		TeamID:     "teamID",
		KeyID:      "keyID",
		KeyContent: []byte{},
		httpClient: mockedHTTPClient,
	}
)

func TestValidateRequest(t *testing.T) {
	var form url.Values
	tokenResponse := TokenResponse{}
	tokenResponseBody, _ := json.Marshal(tokenResponse)
	mockedHTTPClient.On("PostForm", validationEndpoint, form).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(tokenResponseBody)),
		},
		nil,
	)

	res, err := auth.validateRequest(form)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}
