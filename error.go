package apple

import "fmt"

var (
	// ErrorResponseTypeInvalidRequest the request is malformed, typically
	// because it is missing a parameter, contains an unsupported parameter,
	// includes multiple credentials, or uses more than one mechanism for
	// authenticating the client.
	ErrorResponseTypeInvalidRequest ErrorResponseType = "invalid_request"

	// ErrorResponseTypeInvalidClient the client authentication failed,
	// typically due to a mismatched or invalid client identifier, invalid
	// client secret (expired token, malformed claims, or invalid signature), or
	// mismatched or invalid redirect URI.
	ErrorResponseTypeInvalidClient ErrorResponseType = "invalid_client"

	// ErrorResponseTypeInvalidGrant the authorization grant or refresh token
	// is invalid, typically due to a mismatched or invalid client identifier,
	// invalid code (expired or previously used authorization code), or invalid
	// refresh token.
	ErrorResponseTypeInvalidGrant ErrorResponseType = "invalid_grant"

	// ErrorResponseTypeUnauthorizedClient the client is not authorized to use
	// this authorization grant type.
	ErrorResponseTypeUnauthorizedClient ErrorResponseType = "unauthorized_client"

	// ErrorResponseTypeUnsupportedGrantType the authenticated client is not
	// authorized to use this grant type.
	ErrorResponseTypeUnsupportedGrantType ErrorResponseType = "unsupported_grant_type"

	// ErrorResponseTypeInvalidScope the requested scope is invalid.
	ErrorResponseTypeInvalidScope ErrorResponseType = "invalid_scope"

	// ErrorResponseInvalidRequest error when the response is invalid_request. Check message.
	ErrorResponseInvalidRequest = ErrorResponse{
		Type:    ErrorResponseTypeInvalidRequest,
		Message: "The request is malformed, typically because it is missing a parameter, contains an unsupported parameter, includes multiple credentials, or uses more than one mechanism for authenticating the client.",
	}

	// ErrorResponseInvalidClient error when the response is invalid_client. Check message.
	ErrorResponseInvalidClient = ErrorResponse{
		Type:    ErrorResponseTypeInvalidClient,
		Message: "The client authentication failed, typically due to a mismatched or invalid client identifier, invalid client secret (expired token, malformed claims, or invalid signature), or mismatched or invalid redirect URI.",
	}

	// ErrorResponseInvalidGrant error when the response is invalid_grant. Check message.
	ErrorResponseInvalidGrant = ErrorResponse{
		Type:    ErrorResponseTypeInvalidGrant,
		Message: "The authorization grant or refresh token is invalid, typically due to a mismatched or invalid client identifier, invalid code (expired or previously used authorization code), or invalid refresh token.",
	}

	// ErrorResponseUnauthorizedClient error when the response is unauthorized_client. Check message.
	ErrorResponseUnauthorizedClient = ErrorResponse{
		Type:    ErrorResponseTypeUnauthorizedClient,
		Message: "The client is not authorized to use this authorization grant type.",
	}

	// ErrorResponseUnsupportedGrantType error when the response is unsupported_grant_type. Check message.
	ErrorResponseUnsupportedGrantType = ErrorResponse{
		Type:    ErrorResponseTypeUnsupportedGrantType,
		Message: "the authenticated client is not authorized to use this grant type.",
	}

	// ErrorResponseInvalidScope error when the response is invalid_scope. Check message.
	ErrorResponseInvalidScope = ErrorResponse{
		Type:    ErrorResponseTypeInvalidScope,
		Message: "The requested scope is invalid.",
	}
)

type ErrorResponseType string

type ErrorResponse struct {
	Type    ErrorResponseType
	Message string
}

// Error implements the error interface.
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
