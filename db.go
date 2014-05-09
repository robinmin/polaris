package polaris

import (
	"database/sql"
	// "fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-adodb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robinmin/gorp"
	log "github.com/robinmin/logo"
	_ "github.com/ziutek/mymysql/godrv"
	stdlog "log"
	"os"
)

type DBConfig struct {
	// Type is the type of database
	Type string
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
	// _ready is a internal variable to indecate the database is ready or not
	_ready bool
}

type DBEngine struct {
	gorp.DbMap
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
func (conf *DBConfig) InitDB() *DBEngine {
	// get driver
	dialect, driver := dialectAndDriver(conf.Type)
	conf.Driver = driver

	strDSN := conf.String()
	if strDSN == "" {
		log.Error("Invalid DSN has been provided : " + strDSN)
		return nil
	}
	// log.Debug("[" + driver + "] ==> " + strDSN)
	// open connection
	db, err := sql.Open(driver, strDSN)
	if err != nil {
		log.Error("Error connecting to db: " + err.Error())
		return nil
	}
	// fmt.Printf("db = %#v\n", db)
	// dbEngine := &gorp.DbMap{Db: db, Dialect: dialect}
	dbEngine := &DBEngine{DbMap: gorp.DbMap{Db: db, Dialect: dialect}}
	// open SQL log or not
	if conf.Verbose {
		if len(conf.LogFile) > 0 {
			sql_log, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err == nil {
				dbEngine.TraceOn("[SQL]", stdlog.New(sql_log, "sql", stdlog.Lmicroseconds))
			} else {
				dbEngine.TraceOn("[SQL]", log.GetLogger("file"))
			}
		} else {
			dbEngine.TraceOn("[SQL]", log.GetLogger("file"))
		}
	} else {
		dbEngine.TraceOff()
	}

	conf._ready = true
	return dbEngine
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
	return nil, ""
}

//get handler for Martini.Use function
func (conf *DBConfig) MartiniHandler() martini.Handler {
	if conf._ready {
		return func(c martini.Context) {
			dbEngine := conf.InitDB()
			c.Map(dbEngine)
			defer dbEngine.Db.Close()
			c.Next()
		}
	} else {
		// skip to retrieve access in case of not ready
		return func(c martini.Context) { c.Next() }
	}
}

func (egn *DBEngine) Close() bool {
	if egn != nil {
		err := egn.Db.Close()
		if err != nil {
			log.Error("Failed to close DB connecttion : " + err.Error())
			return false
		}
		return true
	}
	return true
}
