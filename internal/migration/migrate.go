package migration

import (
	"os"

	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dir string) error {
	db, err := migrate.New(dir, os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return errors.Wrap(err, "migrate.New")
	}
	err = db.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "db.Up")
	}
	serr, err := db.Close()
	if serr != nil {
		return errors.Wrap(serr, "db.Close source error")
	}
	if err != nil {
		return errors.Wrap(err, "db.Close database error")
	}
	return nil
}

/*
Installing Migration CLI Tool (golang-migrate)
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz
# example: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar -xvz

mv migrate ~/go/bin/
# you can find new releases at: https://github.com/golang-migrate/migrate/releases
*/

/*
To create a migration use the following command:
make migration-create name=<MIGRATION_NAME> dir=<MIGRATION_DIRECTORY>
# example: make migration-create name=partners_init dir=migrations/partners
*/
