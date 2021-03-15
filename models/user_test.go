package models

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/k88t76/CodeArchives-server/config"
)

func TestGetUser(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	expect := &User{
		ID:        1,
		UUID:      "uuid",
		Name:      "name",
		Password:  Encrypt("password"),
		CreatedAt: "2020-01-01",
	}

	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		expect.ID, expect.UUID, expect.Name, expect.Password, expect.CreatedAt)

	user, err := GetUser("uuid")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(user, expect) {
		t.Errorf("user must be %v but %v", expect, user)
	}

}

func TestUser_Validate(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "valid",
			user: User{
				ID:        2,
				UUID:      "uuid2",
				Name:      "name",
				Password:  "password",
				CreatedAt: "2020-01-01",
			},
			want: true,
		},
		{
			name: "invalid (already exist)",
			user: User{
				ID:        3,
				UUID:      "uuid3",
				Name:      "existed name",
				Password:  "password3",
				CreatedAt: "2020-01-10",
			},
			want: false,
		},
	}

	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		3, "uuid3", "existed name", "password3", "2020-01-02")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.Validate(); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCheckUser(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config.DbAccess+"code_archives_test?parseTime=true&loc=Asia%2FTokyo")
	defer func() {
		// DBのCleanup
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()
	tests := []struct {
		name  string
		args  User
		want  string
		want1 bool
		want2 bool
	}{
		{
			name: "valid",
			args: User{
				ID:        4,
				UUID:      "uuid4",
				Name:      "name",
				Password:  "password",
				CreatedAt: "2020-01-01",
			},
			want:  "uuid4",
			want1: true,
			want2: true,
		},
		{
			name: "invalid (unknown username)",
			args: User{
				ID:        5,
				UUID:      "uuid5",
				Name:      "unknown",
				Password:  "password2",
				CreatedAt: "2020-01-02",
			},
			want:  "",
			want1: false,
			want2: false,
		},
		{
			name: "invalid (wrong password)",
			args: User{
				ID:        6,
				UUID:      "uuid6",
				Name:      "name3",
				Password:  "wrong",
				CreatedAt: "2020-01-03",
			},
			want:  "",
			want1: true,
			want2: false,
		},
	}

	// valid user
	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		4, "uuid4", "name", Encrypt("password"), "2020-01-01")

	// wrong password user
	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		6, "uuid6", "name3", Encrypt("incorrect"), "2020-01-03")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := CheckUser(tt.args)

			if got != tt.want {
				t.Errorf("CheckUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckUser() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("CheckUser() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestGetUserNameByToken(t *testing.T) {
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
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			args:    "token",
			want:    "name",
			wantErr: false,
		},
		{
			name:    "invalid (wrong token)",
			args:    "token2",
			want:    "",
			wantErr: true,
		},
	}

	// valid user
	db.MustExec("INSERT INTO users(id, uuid, name, password, created_at) VALUES (?, ?, ?, ?, ?)",
		4, "uuid4", "name", Encrypt("password"), "2020-01-01")
	// valid session
	db.MustExec("INSERT INTO sessions(id, uuid, token, user_id, created_at) VALUES (?, ?, ?, ?, ?)",
		1, "uuid", "token", "uuid4", "2020-01-01")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserNameByToken(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserNameByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserNameByToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
