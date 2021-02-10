package models

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetArchive(t *testing.T) {

	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want *Archive
	}{
		{"test1", args{"uuid"}, &Archive{1, "uuid", "content", "title", "author", "language", "2020-01-01"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := fmt.Sprintf("INSERT INTO %s (id, uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)", tableNameArchives)
			_, err := db.Exec(cmd, tt.want.ID, tt.want.UUID, tt.want.Content, tt.want.Title, tt.want.Author, tt.want.Language, tt.want.CreatedAt)
			if err != nil {
				log.Fatalln(err)
			}
			if got := GetArchive(tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArchive() = %v, want %v", got, tt.want)
			}

			cmd = fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableNameArchives)
			db.Exec(cmd, tt.want.UUID)
		})
	}

}

func TestArchive_Create(t *testing.T) {
	type fields struct {
		ID        int64
		UUID      string
		Content   string
		Title     string
		Author    string
		Language  string
		CreatedAt string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Archive
	}{
		{"test1",
			fields{Content: "content", Title: "title", Author: "author", Language: "language"}, &Archive{Content: "content", Title: "title", Author: "author", Language: "language"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Archive{
				Content:  tt.fields.Content,
				Title:    tt.fields.Title,
				Author:   tt.fields.Author,
				Language: tt.fields.Language,
			}
			a.Create()
			cmd := fmt.Sprintf("SELECT content, title, author, language FROM %s WHERE content = ?", tableNameArchives)
			row := db.QueryRow(cmd, a.Content)
			var archive Archive
			row.Scan(&archive.Content, &archive.Title, &archive.Author, &archive.Language)
			got := NewArchive(archive.ID, "", archive.Content, archive.Title, archive.Author, archive.Language, archive.CreatedAt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}

			cmd = fmt.Sprintf("DELETE FROM %s WHERE content = ?", tableNameArchives)
			db.Exec(cmd, a.Content)
		})
	}
}
