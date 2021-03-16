package archive

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

func TestGetAll(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table archives")
		db.MustExec("truncate table users")
		db.MustExec("truncate table sessions")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := &[]models.Archive{
		{
			ID:        100,
			UUID:      "uuid",
			Content:   "content",
			Title:     "title",
			Author:    "author",
			Language:  "language",
			CreatedAt: "2020-01-03",
		},
		{
			ID:        200,
			UUID:      "uuid2",
			Content:   "content2",
			Title:     "title2",
			Author:    "author",
			Language:  "language2",
			CreatedAt: "2020-01-02",
		},
		{
			ID:        300,
			UUID:      "uuid3",
			Content:   "content3",
			Title:     "title3",
			Author:    "author",
			Language:  "language3",
			CreatedAt: "2020-01-01",
		},
	}

	dummy := &[]models.Archive{
		{
			ID:        400,
			UUID:      "uuid",
			Content:   "content",
			Title:     "title",
			Author:    "dummy",
			Language:  "language",
			CreatedAt: "2020-01-01",
		},
		{
			ID:        500,
			UUID:      "uuid2",
			Content:   "content2",
			Title:     "title2",
			Author:    "dummy",
			Language:  "language2",
			CreatedAt: "2020-01-02",
		},
	}

	for _, a := range *expect {
		db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			a.ID, a.UUID, a.Content, a.Title, a.Author, a.Language, a.CreatedAt)
	}

	for _, a := range *dummy {
		db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			a.ID, a.UUID, a.Content, a.Title, a.Author, a.Language, a.CreatedAt)
	}

	user := models.User{
		ID:        101,
		UUID:      "userid",
		Name:      "author",
		Password:  "password",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.UUID, user.Name, models.Encrypt(user.Password), user.CreatedAt)

	session := models.Session{
		ID:        102,
		UUID:      "sessionid",
		Token:     "token",
		UserID:    "userid",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO sessions(id, uuid, token, user_id, created_at) VALUES (?, ?, ?, ?, ?)",
		session.ID, session.UUID, session.Token, session.UserID, session.CreatedAt)

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&session.Token); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/archives", &body)
	rec := httptest.NewRecorder()
	GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code must be 200 but: %d", rec.Code)
	}

	var archives []models.Archive
	if err := json.NewDecoder(rec.Body).Decode(&archives); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(archives, *expect) {
		t.Errorf("archives must be %v but %v", *expect, archives)
	}
}
