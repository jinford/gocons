package example

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Person struct {
	id        string         `cons:"getter"`
	name      string         `cons:"getter"`
	tags      []string       `cons:"getter"`
	desc      sql.NullString `cons:"getter"`
	createdAt *time.Time     `cons:"getter"`
	*deposit  `cons:"getter"`
}

type deposit struct {
	charge []decimal.Decimal `cons:"getter"`
}
