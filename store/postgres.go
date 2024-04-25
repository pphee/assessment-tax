package store

import (
	"github.com/pphee/assessment-tax/store/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresStore struct {
	DB *gorm.DB
}

func NewPostgresStore(dsn string) *PostgresStore {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting database connection: ", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	db.Exec("CREATE TYPE allowance_type AS ENUM ('Personal', 'Kreceipt');")

	if err := db.AutoMigrate(&modelgorm.AllowanceGorm{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	if err := modelgorm.InitializeData(db); err != nil {
		log.Fatal("Failed to initialize data: ", err)
	}

	return &PostgresStore{DB: db}
}
