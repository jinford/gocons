package example

import "database/sql"

type Person struct {
	id   string         `cons:"getter"`
	name string         `cons:"getter"`
	tags []string       `cons:"getter"`
	desc sql.NullString `cons:"getter"`
}
