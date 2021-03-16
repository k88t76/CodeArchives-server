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

func Test_getArchive(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table archives")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := models.Archive{
		ID:        100,
		UUID:      "uuid",
		Content:   "content",
		Title:     "title",
		Author:    "author",
		Language:  "language",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		expect.ID, expect.UUID, expect.Content, expect.Title, expect.Author, expect.Language, expect.CreatedAt)

	var body bytes.Buffer
	req := httptest.NewRequest(http.MethodGet, "/archive/uuid", &body)
	rec := httptest.NewRecorder()
	getArchive(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code must be 200 but: %d", rec.Code)
	}

	var archive models.Archive
	if err := json.NewDecoder(rec.Body).Decode(&archive); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(archive, expect) {
		t.Errorf("archive must be %v but %v", expect, archive)
	}
}

func Test_postArchive(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table archives")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := models.Archive{
		Content:  "content",
		Title:    "title",
		Author:   "author",
		Language: "language",
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&expect); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/archive/", &body)
	rec := httptest.NewRecorder()
	postArchive(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status code must be 201 but: %d", rec.Code)
	}

	var uuid string
	if err := json.NewDecoder(rec.Body).Decode(&uuid); err != nil {
		t.Fatal(err)
	}

	var archive models.Archive
	if err := db.Get(&archive, "SELECT content, title, author, language FROM archives WHERE uuid = ?", uuid); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(archive, expect) {
		t.Errorf("archive must be %v but %v", expect, archive)
	}
}
