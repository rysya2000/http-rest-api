package apiserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/rysya2000/http-rest-api/internal/app/model"
	"github.com/rysya2000/http-rest-api/internal/app/store/teststore"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	s := newServer(store, *client)
	mw := s.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			// s.sessionStore.Set("user_id", u.ID, 30*time.Minute)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, uuid.NewV4().String()))
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

// func TestServer_HandleUsersCreate(t *testing.T) {
// 	mr, err := miniredis.Run()
// 	if err != nil {
// 		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	client := redis.NewClient(&redis.Options{
// 		Addr: mr.Addr(),
// 	})
// 	s := newServer(teststore.New(), *client)
// 	testCases := []struct {
// 		name         string
// 		payload      interface{}
// 		expectedCode int
// 	}{
// 		{
// 			name: "valid",
// 			payload: map[string]string{
// 				"email":    "user@example.org",
// 				"password": "password",
// 			},
// 			expectedCode: http.StatusCreated,
// 		},
// 		{
// 			name:         "invalid payload",
// 			payload:      "invalid",
// 			expectedCode: http.StatusBadRequest,
// 		},
// 		{
// 			name: "invalid params",
// 			payload: map[string]string{
// 				"email": "invalid",
// 			},
// 			expectedCode: http.StatusUnprocessableEntity,
// 		},
// 	}
// }

// func TestServer_HandleSessionsCreate(t *testing.T) {
// 	u := model.TestUser(t)
// 	store := teststore.New()
// 	store.User().Create(u)
// 	mr, err := miniredis.Run()
// 	if err != nil {
// 		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	client := redis.NewClient(&redis.Options{
// 		Addr: mr.Addr(),
// 	})
// 	s := newServer(store, *client)
// 	testCases := []struct {
// 		name         string
// 		payload      interface{}
// 		expectedCode int
// 	}{
// 		{
// 			name: "valid",
// 			payload: map[string]string{
// 				"email":    u.Email,
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusOK,
// 		},
// 		{
// 			name: "invalid payload",
// 			payload: map[string]string{
// 				"email":    "invalid",
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rec := httptest.NewRecorder()
// 			b := &bytes.Buffer{}
// 			json.NewEncoder(b).Encode(tc.payload)
// 			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
// 			s.ServeHTTP(rec, req)
// 			assert.Equal(t, tc.expectedCode, rec.Code)
// 		})
// 	}
// }

// func TestServer_HandleSeesionsCreate(t *testing.T) {
// 	u := model.TestUser(t)
// 	store := teststore.New()
// 	store.User().Create(u)

// 	s := newServer(store)
// 	testCases := []struct {
// 		name         string
// 		payload      interface{}
// 		expectedCode int
// 	}{
// 		{
// 			name: "valid",
// 			payload: map[string]string{
// 				"email":    u.Email,
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusOK,
// 		},
// 		{
// 			name: "invalid payload",
// 			payload: map[string]string{
// 				"email":    "invalid",
// 				"password": u.Password,
// 			},
// 			expectedCode: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rec := httptest.NewRecorder()
// 			b := &bytes.Buffer{}
// 			json.NewEncoder(b).Encode(tc.payload)
// 			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
// 			s.ServeHTTP(rec, req)
// 			assert.Equal(t, tc.expectedCode, rec.Code)
// 		})
// 	}
// }
