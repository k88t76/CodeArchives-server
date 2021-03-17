package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/k88t76/CodeArchives-server/config"
	"github.com/k88t76/CodeArchives-server/models"
)

func TestSignIn(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("truncate table sessions")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	tests := []struct {
		name       string
		args       models.User
		wantToken  string
		wantStatus int
	}{
		{
			name: "valid",
			args: models.User{
				Name:     "name",
				Password: "password",
			},
			wantToken:  "",
			wantStatus: 200,
		},
		{
			name: "invalid (unknown username)",
			args: models.User{
				Name:     "unknown name",
				Password: "password",
			},
			wantToken:  "Unknown User",
			wantStatus: 401,
		},
		{
			name: "invalid (wrong password)",
			args: models.User{
				Name:     "name",
				Password: "wrong password",
			},
			wantToken:  "Wrong Password",
			wantStatus: 401,
		},
	}

	user := models.User{
		UUID:      "userid",
		Name:      "name",
		Password:  "password",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO users(uuid, name, password, created_at) VALUES (?, ?, ?, ?)",
		user.UUID, user.Name, models.Encrypt(user.Password), user.CreatedAt)

	session := models.Session{
		UUID:      "sessionid",
		Token:     "token",
		UserID:    "userid",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO sessions(uuid, token, user_id, created_at) VALUES (?, ?, ?, ?)",
		session.UUID, session.Token, session.UserID, session.CreatedAt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body bytes.Buffer
			if err := json.NewEncoder(&body).Encode(&tt.args); err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/signin", &body)
			rec := httptest.NewRecorder()
			SignIn(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status code must be %d but: %d", tt.wantStatus, rec.Code)
			}

			if tt.wantToken == "" {
				if err := db.Get(&tt.wantToken, "SELECT token FROM sessions WHERE user_id = ?", user.UUID); err != nil {
					t.Fatal(err)
				}
			}

			var token string
			if err := json.NewDecoder(rec.Body).Decode(&token); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("token must be %v but %v", tt.wantToken, token)
			}

		})
	}

}
