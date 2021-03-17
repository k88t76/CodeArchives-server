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

func TestSignUp(t *testing.T) {
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
				Name:     "valid name",
				Password: "password1",
			},
			wantToken:  "",
			wantStatus: 200,
		},
		{
			name: "invalid (used username)",
			args: models.User{
				Name:     "used name",
				Password: "password2",
			},
			wantToken:  "UsedName",
			wantStatus: 401,
		},
	}

	user := models.User{
		UUID:      "userid",
		Name:      "used name",
		Password:  "password2",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO users(uuid, name, password, created_at) VALUES (?, ?, ?, ?)",
		user.UUID, user.Name, models.Encrypt(user.Password), user.CreatedAt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body bytes.Buffer
			if err := json.NewEncoder(&body).Encode(&tt.args); err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/signup", &body)
			rec := httptest.NewRecorder()
			SignUp(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status code must be %d but: %d", tt.wantStatus, rec.Code)
			}

			if tt.wantToken == "" {
				var name string
				if err := db.Get(&name, "SELECT name FROM users WHERE password = ?", models.Encrypt(tt.args.Password)); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(name, tt.args.Name) {
					t.Errorf("name must be %v but %v", tt.args.Name, name)
				}

			} else {

				var token string
				if err := json.NewDecoder(rec.Body).Decode(&token); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(token, tt.wantToken) {
					t.Errorf("token must be %v but %v", tt.wantToken, token)
				}
			}

		})
	}
}
