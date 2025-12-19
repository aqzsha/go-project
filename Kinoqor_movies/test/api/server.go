package api

import (
	"context"
	"movies/api"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"gorm.io/gorm"
)

type TestServer struct {
	t      *testing.T
	expect *httpexpect.Expect
	db     *gorm.DB
}

func NewTestServer(t *testing.T) *TestServer {
	t.Helper()

	ctx := context.Background()

	engine, db, err := api.NewTestEngine(ctx)
	if err != nil {
		t.Fatalf("failed to init test engine: %v", err)
	}

	// truncateTables(ctx, db)

	ts := httptest.NewServer(engine)
	t.Cleanup(ts.Close)

	config := httpexpect.Config{
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
		BaseURL: ts.URL,
	}

	e := httpexpect.WithConfig(config)

	return &TestServer{
		t:      t,
		expect: e,
		db:     db,
	}
}

func (ts *TestServer) DB() *gorm.DB {
	return ts.db
}

func (ts *TestServer) Expect() *httpexpect.Expect {
	return ts.expect
}

func truncateTables(ctx context.Context, db *gorm.DB) {
	var tables []string

	db.Raw(`
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public'
          AND tablename NOT IN ('goose_db_version')
    `).Scan(&tables)

	for _, table := range tables {
		db.WithContext(ctx).Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
	}
}