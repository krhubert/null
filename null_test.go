package null

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/shopspring/decimal"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "14", []string{"POSTGRES_HOST_AUTH_METHOD=trust"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		conn := fmt.Sprintf("postgresql://postgres:postgres@localhost:%s/postgres", resource.GetPort("3306/tcp"))
		cfg, err := pgxpool.ParseConfig(conn)
		if err != nil {
			return err
		}

		db, err = pgxpool.ConnectConfig(context.Background(), cfg)
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestDecimal(t *testing.T) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "create table if not exists tdec ( n int );")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(ctx, "insert into tdec values ($1)", New[decimal.Decimal](decimal.NewFromInt(8)))
	if err != nil {
		t.Fatal(err)
	}

	var n Null[decimal.Decimal]
	if err := db.QueryRow(ctx, "select * from tdec limit 1").Scan(&n); err != nil {
		t.Fatal(err)
	}

	if n.State != Set || !n.V.Equal(decimal.NewFromInt(8)) {
		t.Fatal("not ok", n)
	}
}

func TestInt8(t *testing.T) {
	ctx := context.Background()
	_, err := db.Exec(ctx, "create table if not exists tint8 ( n int );")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(ctx, "insert into tint8 values ($1)", New[int8](8))
	if err != nil {
		t.Fatal(err)
	}

	var n Null[int8]
	if err := db.QueryRow(ctx, "select * from tint8 limit 1").Scan(&n); err != nil {
		t.Fatal(err)
	}

	if n.State != Set || n.V != 8 {
		t.Fatal("not ok", n)
	}
}
