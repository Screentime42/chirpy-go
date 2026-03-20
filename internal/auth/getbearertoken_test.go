package auth_test

import (
    "net/http"
    "testing"

    "github.com/Screentime42/chirpy-go/internal/auth"
)

func TestGetBearerToken(t *testing.T) {
    tests := []struct {
        name    string
        header  http.Header
        want    string
        errMsg  string
    }{
        {
            name: "valid token",
            header: http.Header{
                "Authorization": []string{"Bearer abc123"},
            },
            want: "abc123",
        },
        {
            name:   "missing header",
            header: http.Header{},
            errMsg: "authorization header missing",
        },
        {
            name: "wrong prefix",
            header: http.Header{
                "Authorization": []string{"Token abc123"},
            },
            errMsg: "invalid authorization header",
        },
        {
            name: "empty token",
            header: http.Header{
                "Authorization": []string{"Bearer   "},
            },
            errMsg: "invalid authorization header",
        },
        {
            name: "extra whitespace",
            header: http.Header{
                "Authorization": []string{"   Bearer    xyz789   "},
            },
            want: "xyz789",
        },
        {
            name: "lowercase bearer",
            header: http.Header{
                "Authorization": []string{"bearer token123"},
            },
            want: "token123",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := auth.GetBearerToken(tt.header)

            // error expectation
            if tt.errMsg != "" {
                if err == nil {
                    t.Fatalf("expected error %q, got nil", tt.errMsg)
                }
                if err.Error() != tt.errMsg {
                    t.Fatalf("expected error %q, got %q", tt.errMsg, err.Error())
                }
                return
            }

            // no error expected
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if got != tt.want {
                t.Fatalf("expected token %q, got %q", tt.want, got)
            }
        })
    }
}
