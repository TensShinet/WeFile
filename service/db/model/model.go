package model

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/db/conn"
	"time"
)

var logger = logging.GetLogger("db_service_model")

// model Init 一定要在 conn Init 之后
func Init() {
	db := conn.GetDB()
	if db == nil {
		logger.Panic("db is not initialized")
		return
	}
	// 修改 charset
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&User{}, &File{}, &UserFile{}, &Session{}); err != nil {
		logger.Panicf("db AutoMigrate failed, for the reason:%v", err.Error())
	}
}

type User struct {
	ID             int64 `gorm:"primary_key:true"`
	RoleID         int64
	Name           string `gorm:"index;size:64"`
	Password       string `gorm:"size:256"`
	Email          string `gorm:"uniqueIndex;size:64"`
	Phone          string `gorm:"size:64"`
	EmailValidated bool
	PhoneValidated bool
	SignUpAt       time.Time
	LastActiveAt   time.Time
	Profile        string `gorm:"size:256"`
	Status         int
}

type UserFile struct {
	ID           int64  `gorm:"primary_key:true"`
	UserID       int64  `gorm:"uniqueIndex:idx_user_file_name"`
	FileID       int64  `gorm:"uniqueIndex:idx_user_file_name"`
	FileName     string `gorm:"uniqueIndex:idx_user_file_name;size:128"`
	IsDirectory  bool
	Directory    string
	UploadAt     time.Time
	LastUpdateAt time.Time
	Status       int
}

type File struct {
	ID            int64  `gorm:"primary_key:true"`
	Hash          string `gorm:"index;size:64"`
	HashAlgorithm string `gorm:"size:64"`
	Size          int64
	Location      string
	CreateAt      time.Time
	UpdateAt      time.Time
	Status        int
}

type Session struct {
	ID        int64  `gorm:"primary_key:true"`
	Token     string `gorm:"size:32"`
	UserID    int64  `gorm:"uniqueIndex;"`
	CreateAt  time.Time
	ExpireAt  time.Time
	CSRFToken string `gorm:"size:32"`
}
