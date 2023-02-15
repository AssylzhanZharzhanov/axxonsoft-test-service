package postgres

import (
	"fmt"
	"log"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConnection creates a new database connection.
func NewConnection(dsn string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to setup sql database: %s", err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// Paginate - gorm pagination.
func Paginate(page domain.PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Offset(page.Offset * page.Size).Limit(page.Size)
		return db
	}
}
