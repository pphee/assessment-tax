package modelgorm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestInitializeData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	mock.ExpectBegin()

	for _, cfg := range []struct {
		AllowanceType string
		Amount        float64
	}{
		{"PersonalDefault", 60000},
		{"PersonalMax", 100000},
		{"DonationMax", 100000},
		{"KReceiptDefault", 50000},
		{"KReceiptMax", 100000},
	} {
		sqlQuery := `SELECT \* FROM "allowance_gorms" WHERE "allowance_gorms"\."allowance_type" = \$1 ORDER BY "allowance_gorms"\."id" LIMIT \$2`
		mock.ExpectQuery(sqlQuery).
			WithArgs(cfg.AllowanceType, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "allowance_type", "amount"}).AddRow(1, cfg.AllowanceType, cfg.Amount))
	}

	mock.ExpectCommit()

	if err := InitializeData(gormDB); err != nil {
		t.Errorf("InitializeData failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
