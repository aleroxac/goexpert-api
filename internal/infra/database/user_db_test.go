package database

import (
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/aleroxac/goexpert-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserDB_CreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "j@j.com", "123456")
	userDB := NewUserDB(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.Find(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}

func TestUserDB_FindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, err := entity.NewUser("John", "j@j.com", "123456")
	assert.Nil(t, err)

	userDB := NewUserDB(db)
	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, user.Password)
}
