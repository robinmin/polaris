// PolarisConfig defined the system configuration and router mapping. It's an interface to the end user to enable them to extend the capability.
//
package polaris

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

// // NewConfig helps to create config object
// func NewConfig(strStatic string, strTemp string, strLog string, strExtension string, strEncoding string, strSessStore string, strSessName string, strSessMask string, redis *RedisConfig, strDBType string, objDB *DBConfig) *PolarisConfig {
// 	return &PolarisConfig{
// 		DirStatic:     strStatic,
// 		DirTemplate:   strTemp,
// 		DirLog:        strLog,
// 		TempExtension: strExtension,
// 		TempEncoding:  strEncoding,
// 		SessionStore:  strSessStore,
// 		SessionName:   strSessName,
// 		SessionMask:   strSessMask,
// 		Redis:         redis,
// 		DBType:        strDBType,
// 		Database:      objDB,
// 	}
// }

// RoterMap provides a interface to the end user to empower them to customize their own router
func (config *PolarisConfig) RoterMap(app *PolarisApplication) bool {
	app.Get("/", func() string {
		return "Hello, you should define your own configuration as the first step!"
	})
	return true
}
