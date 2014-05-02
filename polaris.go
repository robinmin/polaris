// Package polaris provides the major features to the end user. So far, it's a customized version of https://github.com/go-martini/martini/ with other packages.
// It proovides a solid project start in golang.
//
package polaris

import (
	// "database/sql"
	// _ "github.com/mattn/go-adodb"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/boj/redistore"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"github.com/robinmin/gorp"
	log "github.com/robinmin/logo"
	stdlog "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// var (
// 	log_file   string
// 	log_handle *os.File
// )

// PolarisApplication represents a Martini with some reasonable defaults same as the original ClassicMartini.
type PolarisApplication struct {
	*martini.Martini
	martini.Router
	// Config is an user provided object to store relevant information
	Config *PolarisConfig
	// CfgFile is a reference pointer to the config file
	CfgFile *goconfig.ConfigFile
	// Store is a internal variable to store the RediStore object
	Store *redistore.RediStore
	// DbEngine is the pointer to the DB query engine
	DbEngine *gorp.DbMap
	// LogFile is the file name of current log
	LogFile string
	// LogHandle is the log handl
	LogHandle *os.File
}

// NewApp creates a application object with some basic default middleware. It's based on ClassicMartini.
// Classic also maps martini.Routes as a service.
func NewApp(config *PolarisConfig) *PolarisApplication {
	// Add file log
	var log_file string
	var log_handle *os.File
	var err error
	if len(config.DirLog) > 0 {
		year, month, day := time.Now().Date()
		log_file = config.DirLog + "/" + fmt.Sprintf("app_%04d%02d%02d.log", year, month, day)

		log_handle, err = os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Error("error opening file: %v", err)
			return nil
		}
		log.AddLogger("file", log_handle, log.ALL)
	}

	log.Debug("Application started......")
	// create router object
	rtr := martini.NewRouter()

	// create matini object
	server := martini.New()

	// replace log for testing
	server.Map((*log.LevelBasedLogger)(nil))
	server.Use(Logger())

	server.Use(martini.Recovery())
	server.Use(martini.Static(config.DirStatic))
	server.MapTo(rtr, (*martini.Routes)(nil))
	server.Action(rtr.Handle)

	return &PolarisApplication{server, rtr, config, nil, nil, nil, log_file, log_handle}
}

// NewAppWithConfig creates a application object with specified config file.
func NewAppWithConfig(strConfig string) *PolarisApplication {
	cfg, err := goconfig.LoadConfigFile(strConfig)
	if err != nil {
		log.Error("Failed to load config file：", strConfig, ":", err)
		return nil
	}
	dbConfig := &DBConfig{
		Driver:   cfg.MustValue("database", "db_type", "mssql"),
		Host:     cfg.MustValue("database", "db_server", "localhost"),
		Port:     cfg.MustValue("database", "db_port", "1433"),
		Database: cfg.MustValue("database", "db_database", "temdb"),
		User:     cfg.MustValue("database", "db_user", "sa"),
		Password: cfg.MustValue("database", "db_password", ""),
		Verbose:  cfg.MustBool("database", "db_verbose", true),
		LogFile:  cfg.MustValue("database", "db_log", ""),
	}
	redisCOnfig := &RedisConfig{
		Size:     cfg.MustInt("redis", "redis_maxidel", 10),
		Network:  cfg.MustValue("redis", "redis_network", "tcp"),
		Address:  cfg.MustValue("redis", "redis_address", "localhost:6379"),
		Password: cfg.MustValue("redis", "redis_password", ""),
		DB:       cfg.MustValue("redis", "redis_DB", "0"),
	}

	cfgItem := &PolarisConfig{
		DirStatic:     cfg.MustValue("system", "dir_public", "public"),
		DirTemplate:   cfg.MustValue("system", "dir_template", "view"),
		DirLog:        cfg.MustValue("system", "dir_log", "."),
		TempExtension: cfg.MustValue("system", "tpl_extension", ".tpl"),
		TempEncoding:  cfg.MustValue("system", "tpl_encoding", "UTF-8"),
		SessionStore:  cfg.MustValue("session", "session_store", "redis"),
		SessionName:   cfg.MustValue("session", "session_name", "my_session"),
		SessionMask:   cfg.MustValue("session", "session_mask", ""),
		Redis:         redisCOnfig,
		DBType:        dbConfig.Driver,
		Database:      dbConfig,
		url: map[string]string{
			"RedirectUrl":   cfg.MustValue("system", "RedirectUrl", "/new-login"),
			"RedirectParam": cfg.MustValue("system", "RedirectParam", "new-next"),
		},
		Port: cfg.MustInt("system", "port", 3000),
	}
	// create the app server
	app := NewApp(cfgItem)
	app.CfgFile = cfg
	return app
}

func (app *PolarisApplication) RotateLog() bool {
	// get new log
	year, month, day := time.Now().Date()
	new_log := app.Config.DirLog + "/" + fmt.Sprintf("app_%04d%02d%02d.log", year, month, day)
	new_handle, err := os.OpenFile(new_log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error("error opening file: %v", err)
		return false
	}

	// replace it
	if app.LogHandle != nil {
		app.LogHandle.Close()
	}
	app.LogHandle = new_handle
	app.LogFile = new_log
	return true
}

func (app *PolarisApplication) Init(newUser func() sessionauth.User) bool {
	log.Debug("Initializing application......")

	// add middleware -- martini-contrib/render
	app.Use(render.Renderer(render.Options{
		Directory:  app.Config.DirTemplate,
		Extensions: []string{app.Config.TempExtension},
		Charset:    app.Config.TempEncoding,
	}))

	// add middleware -- martini-contrib/sessions
	if "redis" == strings.ToLower(strings.TrimSpace(app.Config.SessionStore)) {
		var err error
		log.Debug("Connect to redis......" + app.Config.Redis.Address)
		if app.Store == nil {
			app.Store, err = redistore.NewRediStoreWithDB(
				app.Config.Redis.Size,
				app.Config.Redis.Network,
				app.Config.Redis.Address,
				app.Config.Redis.Password,
				app.Config.Redis.DB,
				[]byte(app.Config.SessionMask),
			)
			if err != nil {
				log.Error("Failed to connect to redis : " + err.Error())
				return false
			}
		}
		app.Use(sessions.Sessions(app.Config.SessionName, app.Store))
	} else {
		app.Use(sessions.Sessions(app.Config.SessionName, sessions.NewCookieStore([]byte(app.Config.SessionMask))))
	}

	// Connect to databse
	if len(app.Config.Database.Database) > 0 {
		app.DbEngine = InitDB(
			app.Config.DBType,
			app.Config.Database.Host,
			app.Config.Database.Port,
			app.Config.Database.Database,
			app.Config.Database.User,
			app.Config.Database.Password,
			app.Config.Database.Verbose,
			app.Config.Database.LogFile,
		)
		if app.DbEngine == nil {
			log.Error("Failed to connect to database (" + app.Config.Database.Database + ")")
			return false
		}

		// open SQL log or not
		if app.Config.Database.Verbose {
			if len(app.Config.Database.LogFile) > 0 {
				sql_log, err := os.OpenFile(app.Config.Database.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Error("error opening file: %v", err)
					return false
				}
				app.DbEngine.TraceOn("[SQL]", stdlog.New(sql_log, "sql", stdlog.Lmicroseconds))
			} else {
				app.DbEngine.TraceOn("[SQL]", log.GetLogger("file"))
			}
		} else {
			app.DbEngine.TraceOff()
		}
	}

	app.Use(sessionauth.SessionUser(newUser))
	sessionauth.RedirectUrl = app.Config.url["RedirectUrl"]
	sessionauth.RedirectParam = app.Config.url["RedirectParam"]

	// call user defined touter map
	return app.Config.RoterMap(app)
}

func (app *PolarisApplication) UnInit() bool {
	log.Debug("Uninitializing application......")
	// close the redis session handle
	if app.Store != nil {
		app.Store.Close()
		app.Store = nil
	}
	// close the handle of log file
	if app.LogHandle != nil {
		app.LogHandle.Close()
		app.LogHandle = nil
	}
	return true
}

func (app *PolarisApplication) RunApp() bool {
	host := os.Getenv("HOST")
	port := app.Config.Port
	log.Info("listening on " + host + ":" + strconv.Itoa(port))
	lgr := log.GetLogger("stdout")
	if lgr != nil {
		lgr.Fatalln(http.ListenAndServe(host+":"+strconv.Itoa(port), app))
		return true
	}
	return false
}

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, lgr *log.LevelBasedLogger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		log.Debug("==>Started %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(martini.ResponseWriter)
		c.Next()

		log.Debug("<==Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
