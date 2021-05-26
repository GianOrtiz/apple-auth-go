# apple-auth-go

[![GoDoc](https://godoc.org/github.com/GianOrtiz/apple-auth-go?status.svg)](https://pkg.go.dev/github.com/GianOrtiz/apple-auth-go)
[![codecov](https://codecov.io/gh/GianOrtiz/apple-auth-go/branch/master/graph/badge.svg)](https://codecov.io/gh/GianOrtiz/apple-auth-go)
[![Build Status](https://travis-ci.com/GianOrtiz/apple-auth-go.svg?branch=master)](https://travis-ci.com/GianOrtiz/apple-auth-go)

`apple-auth-go` is a unofficial Golang package to validate authorization tokens and manage the authorization of Apple Sign In server side. It provides utility functions and models to retrieve user information and validate authorization codes.

## Installation

Install with go modules:

```
go get github.com/GianOrtiz/apple-auth-go
```

## Usage

The package follow the Go approach to resolve problems, the usage is pretty straightforward, you start initiating a client with:

```go
package main

import (
    "github.com/GianOrtiz/apple-auth-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }
}
```

To validate an authorization code, retrieving refresh and access tokens:

```go
package main

import (
    "github.com/GianOrtiz/apple-auth-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from a mobile app.
    tokenResponse, err := appleAuth.ValidateCode("<AUTHORIZATION-CODE>")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from web app with redirect uri.
    tokenResponse, err := appleAuth.ValidateCodeWithRedirectURI("<AUTHORIZATION-CODE>", "https://redirect-uri")
    if err != nil {
        panic(err)
    }
}
```

The returned `tokenResponse` provides the access token, to make requests on behalf of the user with Apple servers, the refresh token, to retrieve a new access token after expiration, trought the `ValidateRefreshToken` method, and the id token, which is a JWT encoded string with user information. To retrieve the user information from this id token we provide a utility function `GetUserInfoFromIDToken`:

```go
package main

import (
    "fmt"

    "github.com/GianOrtiz/apple-auth-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from a mobile app.
    tokenResponse, err := appleAuth.ValidateCode("<AUTHORIZATION-CODE>")
    if err != nil {
        panic(err)
    }

    user, err := apple.GetUserInfoFromIDToken(tokenResponse.idToken)
    if err != nil {
        panic(err)
    }

    // User Apple unique identification.
    fmt.Println(user.UID)
    // User email if the user provided it.
    fmt.Println(user.Email)
}
```
