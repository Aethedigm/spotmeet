package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGracefully(err)
	}

	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	cel.RootPath = path
	cel.DB.DataType = os.Getenv("DATABASE_TYPE")
}

func getDSN() string {
	dbType := cel.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		}
		return dsn
	}

	return "mysql://" + cel.BuildDSN()
}

func initDB() {
	dirpath := "db-data/postgres"
	path := []string{
		"base",
		"global",
		"pg_commit_ts",
		"pg_dynshmem",
		"pg_logical",
		"pg_logical/mappings",
		"pg_logical/snapshots",
		"pg_multixact",
		"pg_notify",
		"pg_replslot",
		"pg_serial",
		"pg_snapshots",
		"pg_stat",
		"pg_stat_tmp",
		"pg_subtrans",
		"pg_tblspc",
		"pg_twophase",
		"pg_wal",
		"pg_xact",
	}

	for _, p := range path {
		err := os.Mkdir(dirpath+"/"+p, 0755)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				log.Println(p + " already exists. Skipping...")
			} else {
				log.Println("Error creating directory: ", p)
			}
		}
	}
}

func showHelp() {
	color.Yellow("Usage: celeritas [command] [options]")
	color.Yellow("Commands:")
	color.Yellow("  init [project name] - create a new project")
	color.Yellow("  version - show the version")
	color.Yellow("  help - show this help")
	color.Yellow("  migrate [up|down] [version] - run migrations")
	color.Yellow("  make [migration|model|handler|auth] [name] - create a new migration, model, or handler")
	color.Yellow("  init-db - initialize the necessary database folders")
}
