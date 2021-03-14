package models

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/k88t76/CodeArchives-server/config"
)

func TestGetArchive(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table archives")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := &Archive{
		ID:        1000,
		UUID:      "uuid",
		Content:   "content",
		Title:     "title",
		Author:    "author",
		Language:  "language",
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		expect.ID, expect.UUID, expect.Content, expect.Title, expect.Author, expect.Language, expect.CreatedAt)

	a, err := GetArchive("uuid")
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != expect.ID {
		t.Errorf("id must be %d but %d", expect.ID, a.ID)
	}
	if a.UUID != expect.UUID {
		t.Errorf("uuid must be %s but %s", expect.UUID, a.UUID)
	}
	if a.Content != expect.Content {
		t.Errorf("content must be %s but %s", expect.Content, a.Content)
	}
	if a.Title != expect.Title {
		t.Errorf("title must be %s but %s", expect.Title, a.Title)
	}
	if a.Author != expect.Author {
		t.Errorf("author must be %s but %s", expect.Author, a.Author)
	}
	if a.Language != expect.Language {
		t.Errorf("language must be %s but %s", expect.Language, a.Language)
	}
	if a.CreatedAt != expect.CreatedAt {
		t.Errorf("created_at must be %s but %s", expect.CreatedAt, a.CreatedAt)
	}

}

func TestGetArchivesByUser(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table archives")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := &[]Archive{
		{
			ID:        100,
			UUID:      "uuid1",
			Content:   "content1",
			Title:     "title1",
			Author:    "author",
			Language:  "language1",
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

	dummy := &[]Archive{
		{
			ID:        1001,
			UUID:      "uuid11",
			Content:   "content11",
			Title:     "title11",
			Author:    "dummy",
			Language:  "language11",
			CreatedAt: "2020-01-11",
		},
		{
			ID:        2002,
			UUID:      "uuid22",
			Content:   "content22",
			Title:     "title22",
			Author:    "dummy",
			Language:  "language22",
			CreatedAt: "2020-01-22",
		},
	}

	for _, a := range *expect {
		db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			a.ID, a.UUID, a.Content, a.Title, a.Author, a.Language, a.CreatedAt)
	}

	for _, a := range *dummy {
		db.MustExec("INSERT INTO archives(id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			a.ID, a.UUID, a.Content, a.Title, a.Author, a.Language, a.CreatedAt)

		archives, err := GetArchivesByUser("author", 100)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(archives, *expect) {
			t.Errorf("archives must be %v but %v", *expect, archives)
		}
	}
}
