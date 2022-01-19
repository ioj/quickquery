package quickquery

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TestDB struct {
	db *pgxpool.Pool
}

func NewTestDB(dburl string) *TestDB {
	u, err := url.Parse(dburl)
	if err != nil {
		panic(err)
	}

	if !strings.HasSuffix(u.Path, "_test") {
		panic("database name doesn't have _test suffix: " + u.Path)
	}

	s := &TestDB{}

	s.db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *TestDB) TearDown() {
	s.db.Close()
}

func (s *TestDB) QQ(t *testing.T) *QuickQuery {
	return New(s.db, t)
}
