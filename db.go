package polaris

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-adodb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robinmin/gorp"
	log "github.com/robinmin/logo"
	_ "github.com/ziutek/mymysql/godrv"
)

type DBConfig struct {
	// Driver is the driver name of current database, so far, Polaris support : mysql/gomysql/postgres/sqlite/adodb
	Driver string
	// Host is the host name of current database server
	Host string
	// Port is the port number of current database server
	Port string
	// Database is the database name of current database
	Database string
	// User is the user name to access current database
	User string
	// Password is the password of current user
	Password string
	// Verbose is flag to output SQL statement into log or not
	Verbose bool
	// LogFile is the log file name for SQL statement.
	LogFile string
}

// String is a helper to generate connection string for mssql
func (conf *DBConfig) String() string {
	switch conf.Driver {
	case "mymysql":
		return conf.Database + "/" + conf.User + "/" + conf.Password
	case "mysql":
		return conf.User + ":" + conf.Password + "@/" + conf.Database
	case "postgres":
		return "user=" + conf.User + " password=" + conf.Password + " dbname=" + conf.Database + " sslmode=disable"
	case "sqlite3":
		return conf.Database
	case "adodb":
		return "Provider=SQLOLEDB;Initial Catalog=" + conf.Database + ";Data Source=" + conf.Host + "," + conf.Port + ";User Id=" + conf.User + ";Password=" + conf.Password + ";CharacterSet=utf-8;"
	}
	return ""
}

// InitDB initialize a connection to specified database
func InitDB(strDbType string, strHost string, strPort string, strDtbs string, strUser string, strPass string, blVerbose bool, strLog string) *gorp.DbMap {
	// get driver
	dialect, driver := dialectAndDriver(strDbType)
	// get connection string
	dbConfig := DBConfig{
		Driver:   driver,
		Host:     strHost,
		Port:     strPort,
		Database: strDtbs,
		User:     strUser,
		Password: strPass,
		Verbose:  blVerbose,
		LogFile:  strLog,
	}
	strDSN := dbConfig.String()
	if strDSN == "" {
		log.Error("Invalid DSN has been provided : " + strDSN)
		return nil
	}
	// open connection
	db, err := sql.Open(driver, strDSN)
	if err != nil {
		log.Error("Error connecting to db: " + err.Error())
		return nil
	}

	return &gorp.DbMap{Db: db, Dialect: dialect}
}

// dialectAndDriver is a internal function create gorp.Dialect object
func dialectAndDriver(strDbType string) (gorp.Dialect, string) {
	switch strDbType {
	case "mysql":
		return gorp.MySQLDialect{"InnoDB", "UTF8"}, "mymysql"
	case "gomysql":
		return gorp.MySQLDialect{"InnoDB", "UTF8"}, "mysql"
	case "postgres":
		return gorp.PostgresDialect{}, "postgres"
	case "sqlite":
		return gorp.SqliteDialect{}, "sqlite3"
	case "adodb":
		return gorp.AdodbDialect{"UTF8"}, "adodb"
	}
	panic("GORP_TEST_DIALECT env variable is not set or is invalid. Please see README.md")
}
