package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"
)

type Archive struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	CreatedAt string `json:"createdAt"`
}

func NewArchive(id int64, uuid string, content string, title string, author string, language string, createdAt string) *Archive {
	return &Archive{
		id,
		uuid,
		content,
		title,
		author,
		language,
		createdAt,
	}
}

func GetArchive(uuid string) *Archive {
	cmd := fmt.Sprintf("SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE uuid = ?", tableNameArchives)
	row := db.QueryRow(cmd, uuid)
	var archive Archive
	err := row.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
	if err != nil {
		return nil
	}
	archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
	return NewArchive(archive.ID, archive.UUID, archive.Content, archive.Title, archive.Author, archive.Language, archive.CreatedAt)
}

func TestArchives() ([]Archive, error) {
	var archives []Archive
	cmd := fmt.Sprintf(`SELECT id, uuid, content, title, author, language, created_at FROM %s ORDER BY created_at DESC`, tableNameArchives)
	rows, err := db.Query(cmd)
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive Archive
		rows.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
		archives = append(archives, archive)
	}
	return archives, nil

}

func GetArchivesByUser(userName string, limit int) ([]Archive, error) {
	var archives []Archive
	cmd := fmt.Sprintf(`SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE author = ? ORDER BY created_at DESC LIMIT ?`, tableNameArchives)
	rows, err := db.Query(cmd, userName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive Archive
		rows.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
		archives = append(archives, archive)
	}
	return archives, nil

}

func GetMatchArchive(search string, userName string) ([]Archive, error) {
	var archives []Archive
	cmd := fmt.Sprintf(`SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE author = ? AND ( content LIKE `, tableNameArchives)
	words := strings.Fields(search)
	len := len(words)
	if len == 1 {
		cmd += "'%" + words[0] + "%')"
	} else {
		for i := 0; i < len-1; i++ {
			cmd += "'%" + words[0] + "%' OR content LIKE "
		}
		cmd += "'%" + words[len-1] + "%')"
	}
	rows, err := db.Query(cmd, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive Archive
		rows.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
		archives = append(archives, archive)
	}
	return archives, nil

}

func (a *Archive) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?)", tableNameArchives)
	_, err := db.Exec(cmd, CreateUUID(), a.Content, a.Title, a.Author, a.Language, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (a *Archive) Update() error {
	cmd := fmt.Sprintf("UPDATE %s SET content = ?, title = ?, language = ?, created_at = ? WHERE uuid = ?", tableNameArchives)
	_, err := db.Exec(cmd, a.Content, a.Title, a.Language, time.Now().Format(time.RFC3339), a.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Archive) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableNameArchives)
	_, err := db.Exec(cmd, a.UUID)
	if err != nil {
		return err
	}
	return nil
}

func CreateUUID() string {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return uuid
}
