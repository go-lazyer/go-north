package north

import "database/sql"

type DataSource struct {
	Db           *sql.DB
	DriverName   string
	DaoFilePaths []string
}
type Config struct {
	MaxOpenConns int
	MaxIdleConns int
}

func Open(driverName string, dsn string, config *Config) (*DataSource, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &DataSource{
		Db:         db,
		DriverName: driverName,
	}, nil
}
