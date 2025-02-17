package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/andreiz53/cookinator/token"
	"github.com/andreiz53/cookinator/util"
)

func setAuth(t *testing.T, request *http.Request, tokenMaker token.Maker, authType string, email string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)

	header := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(authHeaderKey, header)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setup         func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuth(t, request, tokenMaker, authHeaderTypeBearer, util.RandomEmail(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidAuthHeaderType",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuth(t, request, tokenMaker, "notBearer", util.RandomEmail(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuth(t, request, tokenMaker, authHeaderTypeBearer, util.RandomEmail(), -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "EmptyAuthHeaderType",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuth(t, request, tokenMaker, "", util.RandomEmail(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			path := "/auth"
			server.router.GET(path, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			tc.setup(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
