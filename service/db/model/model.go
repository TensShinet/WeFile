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
	ID             int64     `gorm:"not null;primary_key:true"`
	RoleID         int64     `gorm:"not null;default:100000"`
	Name           string    `gorm:"not null;index;size:64"`
	Password       string    `gorm:"not null;size:64"`
	Email          string    `gorm:"not null;uniqueIndex;size:64"`
	Phone          string    `gorm:"not null;size:64"`
	EmailValidated bool      `gorm:"not null;default false"`
	PhoneValidated bool      `gorm:"not null;default false"`
	SignUpAt       time.Time `gorm:"not null;default NOW()"`
	LastActiveAt   time.Time `gorm:"not null;default NOW()"`
	Profile        string    `gorm:"not null;size:255"`
	Status         int       `gorm:"not null;default 0"`
}

type UserFile struct {
	ID           int64     `gorm:"not null;primary_key:true"`
	UserID       int64     `gorm:"not null;index:idx_user"`
	Directory    string    `gorm:"not null;default /;size:2048"`
	FileName     string    `gorm:"not null;size:255"`
	Hash         string    `gorm:"not null;uniqueIndex;size:64"` // UserID + Directory + FileName 的 hash 保证唯一性
	FileID       int64     `gorm:"not null;default 0;"`
	IsDirectory  bool      `gorm:"not null;default false;"`
	UploadAt     time.Time `gorm:"not null;default NOW()"`
	LastUpdateAt time.Time `gorm:"not null;default NOW()"`
	Status       int       `gorm:"not null;default 0"`
}

type File struct {
	ID            int64     `gorm:"not null;primary_key:true"`
	Hash          string    `gorm:"not null;uniqueIndex;size:64"`
	HashAlgorithm string    `gorm:"not null;size:32"`
	Size          int64     `gorm:"not null;default 0;"`
	Count         int       `gorm:"not null;default 0;"` // 引用计数
	Location      string    `gorm:"not null;default /;size:2048"`
	CreateAt      time.Time `gorm:"not null;default NOW()"`
	UpdateAt      time.Time `gorm:"not null;default NOW()"`
	Status        int       `gorm:"not null;default 0"`
}

// session 表基本不用了
type Session struct {
	ID        int64     `gorm:"not null;primary_key:true"`
	Token     string    `gorm:"not null;size:32"`
	UserID    int64     `gorm:"not null;uniqueIndex;"`
	CreateAt  time.Time `gorm:"not null;default NOW()"`
	ExpireAt  time.Time `gorm:"not null;default NOW()"`
	CSRFToken string    `gorm:"not null;size:32"`
}
