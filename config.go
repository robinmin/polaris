// PolarisConfig defined the system configuration and router mapping. It's an interface to the end user to enable them to extend the capability.
//
package polaris

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	log "github.com/robinmin/logo"
	"os"
	"time"
)

// PolarisConfig is the struct to define configuration items
type PolarisConfig struct {
	// DirStatic is the default folder name where the public resouces(e.g, image, css, js and etc) are putted into
	DirStatic string
	// DirTemplate is the folder name where the template files located in
	DirTemplate string
	// DirLog is the folder name where the log files located in
	DirLog string
	// TempExtension is the extension name of template files
	TempExtension string
	// TempEncoding is the encoding name of the templates files
	TempEncoding string
	// SessionStore is the store type of the session. it can be cookie or redis so far
	SessionStore string
	// SessionName is the default session name for cookie to store session id
	SessionName string
	// SessionMask is a mask code to the cookie of session id
	SessionMask string
	// Redis is the option set for redis connection
	Redis *RedisConfig
	// DBType is the datbase type. So far, it can be mysql/gomysql/postgres/sqlite/adodb
	DBType string
	// Database is the configuration items for database
	Database *DBConfig
	// URL is a URL collection
	url map[string]string
	// Port is the port number current server listening on
	Port int
	// CfgFile is a reference pointer to the config file
	CfgFile string
	// CfgHandle is a reference pointer to the config file
	CfgHandle *goconfig.ConfigFile

	// LogFile is the file name of current log
	LogFile string
	// LogHandle is the log handl
	LogHandle *os.File
}

// RedisConfig is the struct to define configuration items for redis
type RedisConfig struct {
	// Size is the maximum number of idle connections
	Size int
	// Network is the communication protocal of redis
	Network string
	// Address is the acctual adress of the redis server hostname + port
	Address string
	// Password is the default password of the connection
	Password string
	// DB is the default database in redis
	DB string
}

// Config is the interface to cental config object
type Config interface {
	LoadConfig() bool
	GetBasicConfig() *PolarisConfig
	RoterMap(app *PolarisApplication) bool
}

// LoadConfig creates a config object by specified config file.
func (cfg *PolarisConfig) LoadConfig() bool {
	cfgFile, err := goconfig.LoadConfigFile(cfg.CfgFile)
	if err != nil {
		log.Error("Failed to load config file：", cfg.CfgFile, ":", err)
		return false
	}

	cfg.DirStatic = cfgFile.MustValue("system", "dir_public", "public")
	cfg.DirTemplate = cfgFile.MustValue("system", "dir_template", "view")
	cfg.DirLog = cfgFile.MustValue("system", "dir_log", ".")
	cfg.TempExtension = cfgFile.MustValue("system", "tpl_extension", ".tpl")
	cfg.TempEncoding = cfgFile.MustValue("system", "tpl_encoding", "UTF-8")
	cfg.SessionStore = cfgFile.MustValue("session", "session_store", "redis")
	cfg.SessionName = cfgFile.MustValue("session", "session_name", "my_session")
	cfg.SessionMask = cfgFile.MustValue("session", "session_mask", "")
	cfg.url = map[string]string{
		"RedirectUrl":   cfgFile.MustValue("system", "RedirectUrl", "/new-login"),
		"RedirectParam": cfgFile.MustValue("system", "RedirectParam", "new-next"),
	}
	cfg.Port = cfgFile.MustInt("system", "port", 3000)

	cfg.Redis = &RedisConfig{
		Size:     cfgFile.MustInt("redis", "redis_maxidel", 10),
		Network:  cfgFile.MustValue("redis", "redis_network", "tcp"),
		Address:  cfgFile.MustValue("redis", "redis_address", "localhost:6379"),
		Password: cfgFile.MustValue("redis", "redis_password", ""),
		DB:       cfgFile.MustValue("redis", "redis_DB", "0"),
	}

	cfg.Database = &DBConfig{
		Driver:   cfgFile.MustValue("database", "db_type", "mssql"),
		Host:     cfgFile.MustValue("database", "db_server", "localhost"),
		Port:     cfgFile.MustValue("database", "db_port", "1433"),
		Database: cfgFile.MustValue("database", "db_database", "temdb"),
		User:     cfgFile.MustValue("database", "db_user", "sa"),
		Password: cfgFile.MustValue("database", "db_password", ""),
		Verbose:  cfgFile.MustBool("database", "db_verbose", true),
		LogFile:  cfgFile.MustValue("database", "db_log", ""),
	}
	cfg.DBType = cfg.Database.Driver
	cfg.CfgHandle = cfgFile

	// Add file log
	if len(cfg.DirLog) > 0 {
		year, month, day := time.Now().Date()
		cfg.LogFile = cfg.DirLog + "/" + fmt.Sprintf("app_%04d%02d%02d.log", year, month, day)

		cfg.LogHandle, _ = os.OpenFile(cfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if cfg.LogHandle != nil {
			log.AddLogger("file", cfg.LogHandle, log.ALL)
		}
	}

	return true
}

func (cfg *PolarisConfig) GetBasicConfig() *PolarisConfig {
	return cfg
}

// RoterMap provides a interface to the end user to empower them to customize their own router
func (config *PolarisConfig) RoterMap(app *PolarisApplication) bool {
	app.Get("/", func() string {
		return "Hello, you should define your own configuration as the first step!"
	})
	return true
}
