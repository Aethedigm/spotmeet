// go:build integration

// run tests with this command: go test . --tags integration --count=1

package data

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "celeritas_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var models Models
var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")
	os.Setenv("UPPER_DB_LOG", "ERROR")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		// _ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	err = createTables(testDB)
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	models = New(testDB)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func compileMigrations(ext string) (string, error) {
	var res []string
	err := filepath.WalkDir("../migrations", func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		_, file := filepath.Split(s)

		if strings.HasSuffix(file, ext) {
			f, err := os.ReadFile(s)
			if err != nil {
				return err
			}

			res = append(res, string(f))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(res, " "), nil
}

func createTables(db *sql.DB) error {
	stmt, err := compileMigrations(".up.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func FindIDIn(id int, matches []*Match) bool {
	for _, m := range matches {
		if m.ID == id {
			return true
		}
	}
	return false
}

func ArtistInArray(a Artist, arr []*Artist) bool {
	for _, v := range arr {
		if v.ID == a.ID {
			return true
		}
	}
	return false
}

func MessageInArray(m Message, arr []*Message) bool {
	for _, v := range arr {
		if v.ID == m.ID {
			return true
		}
	}
	return false
}

func ProfileInArray(p Profile, arr []*Profile) bool {
	for _, v := range arr {
		if v.ID == p.ID {
			return true
		}
	}
	return false
}

func SettingsInArray(s Settings, arr []*Settings) bool {
	for _, v := range arr {
		if v.ID == s.ID {
			return true
		}
	}
	return false
}

func LikedArtistInArray(l LikedArtist, arr []*LikedArtist) bool {
	for _, v := range arr {
		if v.ID == l.ID {
			return true
		}
	}
	return false
}
