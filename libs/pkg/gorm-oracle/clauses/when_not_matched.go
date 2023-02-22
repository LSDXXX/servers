package clauses

import (
	"gorm.io/gorm/clause"
)

// WhenNotMatched description
type WhenNotMatched struct {
	clause.Values
	Where clause.Where
}

// Name description
// @receiver w
// @return string
func (w WhenNotMatched) Name() string {
	return "WHEN NOT MATCHED"
}

// Build description
// @receiver w
// @param builder
func (w WhenNotMatched) Build(builder clause.Builder) {
	if len(w.Columns) > 0 {
		if len(w.Values.Values) != 1 {
			panic("cannot insert more than one rows due to Oracle SQL language restriction")
		}

		builder.WriteString(" THEN")
		builder.WriteString(" INSERT ")
		w.Build(builder)

		if len(w.Where.Exprs) > 0 {
			builder.WriteString(w.Where.Name())
			builder.WriteByte(' ')
			w.Where.Build(builder)
		}
	}
}
