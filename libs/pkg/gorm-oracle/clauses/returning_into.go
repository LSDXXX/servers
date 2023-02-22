package clauses

import (
	"gorm.io/gorm/clause"
)

// ReturningInto description
type ReturningInto struct {
	Variables []clause.Column
	Into      []*clause.Values
}
