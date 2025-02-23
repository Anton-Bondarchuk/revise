package main

import "flag"

func main() {
	var configPath, migratonsPath, migratonTable string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&migratonsPath, "migrations", "", "path to migrations folder")
	flag.StringVar(&migratonTable, "migration_table", "migrations", "migration table name")
	flag.Parse()

	if configPath == "" {
		panic("config path is empty")
	}

	if migratonsPath == "" {
		panic("migrations path is empty")
	}

	if migratonTable == "" {
		panic("migration table name is empty")
	}

	// TODO: check library golang-migrate or pgx stack for correct migratoin
}