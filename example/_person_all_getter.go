package example

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Person struct {
	id        string
	name      string
	tags      []string
	desc      sql.NullString
	createdAt *time.Time
	*deposit
}

type deposit struct {
	charge []decimal.Decimal
}
