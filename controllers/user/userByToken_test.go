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

func TestUserByToken(t *testing.T) {
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
		args       string
		wantName   string
		wantStatus int
	}{
		{
			name:       "valid token",
			args:       "validToken",
			wantName:   "validUser",
			wantStatus: 200,
		},
		{
			name:       "invalid token",
			args:       "invalidToken",
			wantName:   "Invalid Token",
			wantStatus: 400,
		},
	}

	user := models.User{
		UUID:      "userid",
		Name:      "validUser",
		Password:  "password2",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO users(uuid, name, password, created_at) VALUES (?, ?, ?, ?)",
		user.UUID, user.Name, models.Encrypt(user.Password), user.CreatedAt)

	session := models.Session{
		UUID:      "sessionid",
		Token:     "validToken",
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

			req := httptest.NewRequest(http.MethodPost, "/userbytoken", &body)
			rec := httptest.NewRecorder()
			UserByToken(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status code must be %d but: %d", tt.wantStatus, rec.Code)
			}

			var name string
			if err := json.NewDecoder(rec.Body).Decode(&name); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(name, tt.wantName) {
				t.Errorf("name must be %v but %v", tt.wantName, name)

			}

		})
	}
}
