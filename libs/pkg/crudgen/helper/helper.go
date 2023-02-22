package helper

import (
	"strings"

	"gorm.io/gorm"
)

const (
	SystemMaxId = 1000
)

type Cond struct {
	Cond   bool
	Result func() string
}

func IfClause(conds []Cond) string {
	judge := func(c Cond) string {
		if c.Cond {
			return c.Result()
		}
		return ""
	}

	clauses := make([]string, len(conds))
	for i, cond := range conds {
		clauses[i] = strings.Trim(judge(cond), " ")
	}
	return " " + strings.Join(clauses, " ")
}

func WhereClause(conds []string) string {
	return joinClause(conds, "WHERE", whereValue, " ")
}

func SetClause(conds []string) string {
	return joinClause(conds, "SET", setValue, ",")
}

func joinClause(conds []string, keyword string, deal func(string) string, sep string) string {
	clauses := make([]string, len(conds))
	for i, clause := range conds {
		clauses[i] = deal(clause)
	}

	sql := trimAll(strings.Join(clauses, sep))
	if sql != "" {
		sql = " " + keyword + " " + sql
	}
	return sql
}

func trimAll(input string) string {
	input = strings.Trim(input, " ")
	lowercase := strings.ToLower(input)
	switch {
	case strings.HasPrefix(lowercase, "and "):
		return input[4:]
	case strings.HasPrefix(lowercase, "or "):
		return input[3:]
	case strings.HasPrefix(lowercase, "xor "):
		return input[4:]
	case strings.HasPrefix(lowercase, ","):
		return input[1:]
	default:
		return input
	}
}

// whereValue append a new condition with prefix "AND"
func whereValue(value string) string {
	value = strings.Trim(value, " ")
	lowercase := strings.ToLower(value)
	switch {
	case lowercase == "":
		return ""
	case strings.HasPrefix(lowercase, "and "):
		return value
	case strings.HasPrefix(lowercase, "or "):
		return value
	case strings.HasPrefix(lowercase, "xor "):
		return value
	default:
		return "AND " + value
	}
}

func setValue(value string) string {
	return strings.Trim(value, ", ")
}

type DAO[T any, M any] interface {
	DB() *gorm.DB
	WithDB(*gorm.DB) T
	Table() string

	Page(page, pageSize int, order string, conds ...M) (result []M, count int64, err error)
	Find(conds M) (result []M, err error)
	Take(order string, conds ...M) (result M, err error)
	Count(conds ...M) (count int64, err error)

	Insert(items ...*M) error
	InsertInBatches(items []*M, size int) error
	UpdateOrCreate(update *M, conds M) error
	Updates(updates *M, conds M) (rowsAffected int64, err error)
	FirstOrCreate(insert *M, conds M) (rowsAffected int64, err error)

	Delete(conds M) (rowsAffected int64, err error)
}
