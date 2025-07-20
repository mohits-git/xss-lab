package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestGetAuthHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      http.Header
		expected    string
		expectError bool
	}{
		{
			name: "Valid header",
			header: http.Header{
				"Authorization": []string{"Bearer validtoken"},
			},
			expected:    "validtoken",
			expectError: false,
		},
		{
			name:        "Missing header",
			header:      http.Header{},
			expected:    "",
			expectError: true,
		},
		{name: "Invalid format",
			header: http.Header{
				"Authorization": []string{"InvalidFormat"},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "Empty token",
			header: http.Header{
				"Authorization": []string{"Bearer "},
			},
			expected:    "",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetAuthHeader(tt.header)
			if (err != nil) != tt.expectError {
				t.Errorf("GetAuthHeader() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if token != tt.expected {
				t.Errorf("GetAuthHeader() = %v, expected %v", token, tt.expected)
			}
		})
	}
}

func TestSetAuthHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		token    string
		expected http.Header
	}{
		{
			name: "Set valid token",
			header: http.Header{
				"Authorization": []string{"OldToken"},
			},
			token: "NewToken",
			expected: http.Header{
				"Authorization": []string{"Bearer NewToken"},
			},
		},
		{
			name:   "Set empty token",
			header: http.Header{},
			token:  "",
			expected: http.Header{
				"Authorization": []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAuthHeader(tt.header, tt.token)
			if len(tt.header["Authorization"]) != len(tt.expected["Authorization"]) ||
				tt.header.Get("Authorization") != tt.expected.Get("Authorization") {
				t.Errorf("SetAuthHeader() = %v, expected %v", tt.header, tt.expected)
			}
		})
	}
}

func TestMakeJwt(t *testing.T) {
	tests := []struct {
		name        string
		userId      string
		tokenSecret string
		expiresIn   int64
		expectError bool
	}{
		{
			name:        "Valid JWT",
			userId:      "user123",
			tokenSecret: "secret",
			expiresIn:   3600,
			expectError: false,
		},
		{
			name:        "Empty user ID",
			userId:      "",
			tokenSecret: "secret",
			expiresIn:   3600,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtStr, err := MakeJwt(tt.userId, tt.tokenSecret, time.Duration(tt.expiresIn)*time.Second)
			if (err != nil) != tt.expectError {
				t.Errorf("MakeJwt() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && jwtStr == "" {
				t.Error("MakeJwt() returned empty token")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		tokenSecret string
		expectError bool
		expectedId  string
	}{
		{
			name: "Valid JWT",
			token: (func() string {
				token, _ := MakeJwt("user123", "secret", time.Hour)
				return token
			})(),
			tokenSecret: "secret",
			expectError: false,
			expectedId:  "user123",
		},
		{
			name:        "Invalid JWT",
			token:       "invalidtoken",
			tokenSecret: "secret",
			expectError: true,
			expectedId:  "",
		},
		{
			name:        "Empty token",
			token:       "",
			tokenSecret: "secret",
			expectError: true,
			expectedId:  "",
		},
		{
			name:        "Empty secret",
			token:       "validtoken",
			tokenSecret: "",
			expectError: true,
			expectedId:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, err := ValidateJWT(tt.token, tt.tokenSecret)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateJWT() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if userId != tt.expectedId {
				t.Errorf("ValidateJWT() = %v, expected %v", userId, tt.expectedId)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "mypassword",
			expectError: false,
		},
		{
			name:        "Empty password",
			password:    "",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("HashPassword() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && hashed == "" {
				t.Error("HashPassword() returned empty hash")
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	hashed, _ := HashPassword("mypassword")
	tests := []struct {
		name        string
		password    string
		hashed      string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "mypassword",
			hashed:      hashed,
			expectError: false,
		},
		{
			name:        "Invalid password",
			password:    "wrongpassword",
			hashed:      hashed,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hashed, tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("ComparePassword() error = %v, expectError %v", err, tt.expectError)
				return
			}
		})
	}
}

