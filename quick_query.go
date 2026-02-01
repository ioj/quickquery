package quickquery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DefaultTestTimeout = 10 * time.Second

type QuickQuery struct {
	Ctx context.Context
	A   *assert.Assertions
	R   *require.Assertions
	Tx  pgx.Tx

	cancel context.CancelFunc
}

func New(db *pgxpool.Pool, t *testing.T) *QuickQuery {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTestTimeout)
	tx, err := db.Begin(ctx)

	a := assert.New(t)
	r := require.New(t)

	r.NoError(err, "creating transaction")

	return &QuickQuery{Ctx: ctx, A: a, R: r, Tx: tx, cancel: cancel}
}

func (qq *QuickQuery) Done() {
	qq.Tx.Rollback(qq.Ctx)
	qq.cancel()
}

func (qq *QuickQuery) Bool(sql string, args ...interface{}) bool {
	var x bool
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) Time(sql string, args ...interface{}) time.Time {
	var x time.Time
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) Int(sql string, args ...interface{}) int {
	var x int
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) Int64(sql string, args ...interface{}) int64 {
	var x int64
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) String(sql string, args ...interface{}) string {
	var x string
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) IsNull(sql string, args ...interface{}) {
	var x *string
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	qq.A.Nil(x, fmt.Sprintf("%v is not null", sql))
}

func (qq *QuickQuery) Exec(sql string, args ...interface{}) {
	_, err := qq.Tx.Exec(qq.Ctx, sql, args...)
	qq.A.NoError(err, fmt.Sprintf("%v failed", sql))
}

func (qq *QuickQuery) UUID(sql string, args ...interface{}) uuid.UUID {
	var x uuid.UUID
	qq.A.NoError(qq.Tx.QueryRow(qq.Ctx, sql, args...).Scan(&x))
	return x
}

func (qq *QuickQuery) Throws(sqlstate string, sql string, args ...interface{}) *pgconn.PgError {
	_, err := qq.Tx.Exec(qq.Ctx, sql, args...)
	if err == nil {
		qq.A.Error(err)
		return nil
	}

	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		qq.A.NoError(err)
		return nil
	}

	if sqlstate == "" {
		return pgerr
	}

	qq.A.Equal(sqlstate, pgerr.Code, "pg error code other than expected")

	return pgerr
}

func PGErrCode(err error) string {
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return ""
	}

	return pgerr.Code
}
