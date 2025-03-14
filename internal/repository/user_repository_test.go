package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepositoryCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		model model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				model: model.User{
					Username: "testusername",
					Password: "password",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						args.model.Username,
						args.model.Password,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				model: model.User{
					Username: "testusername",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						args.model.Username,
						args.model.Password,
					).
					WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &userRepository{
				db: gormDB,
			}
			if err := r.CreateUser(&tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("repository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepositoryFindUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := &userRepository{db: gormDB}

	tests := []struct {
		name     string
		username string
		wantUser *model.User
		wantErr  bool
		mockFn   func()
	}{
		{
			name:     "success",
			username: "testusername",
			wantUser: &model.User{
				Model:    gorm.Model{ID: 1}, // ✅ Tambahkan gorm.Model agar ID cocok
				Username: "testusername",
				Password: "hashedpassword",
			},
			wantErr: false,
			mockFn: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, "testusername", "hashedpassword")

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
					WithArgs("testusername", 1).
					WillReturnRows(rows)
			},
		},
		{
			name:     "error - user not found",
			username: "unknownuser",
			wantUser: nil,
			wantErr:  true,
			mockFn: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
					WithArgs("unknownuser", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			user, err := repo.FindUserByUsername(tt.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				assert.Equal(t, tt.wantUser.ID, user.ID) // ✅ Cocokkan ID
				assert.Equal(t, tt.wantUser.Username, user.Username)
				assert.Equal(t, tt.wantUser.Password, user.Password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
