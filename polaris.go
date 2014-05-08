// Package polaris provides the major features to the end user. So far, it's a customized version of https://github.com/go-martini/martini/ with other packages.
// It proovides a solid project start in golang.
//
package polaris

import (
	// "database/sql"
	// _ "github.com/mattn/go-adodb"
	"fmt"
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
	Config Config
	// Store is a internal variable to store the RediStore object
	Store *redistore.RediStore
	// DbEngine is the pointer to the DB query engine
	DbEngine *gorp.DbMap
}

// NewApp creates a application object with some basic default middleware. It's based on ClassicMartini.
// Classic also maps martini.Routes as a service.
func NewApp(cfg Config, newUser func() sessionauth.User) *PolarisApplication {
	// Create config object
	if !cfg.LoadConfig() {
		log.Error("Failed to init config object")
		return nil
	}
	config := cfg.GetBasicConfig()

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

	app := &PolarisApplication{server, rtr, cfg, nil, nil}

	log.Debug("Add middleware -- martini-contrib/render......")
	app.Use(render.Renderer(render.Options{
		Directory:  config.DirTemplate,
		Extensions: []string{config.TempExtension},
		Charset:    config.TempEncoding,
	}))

	log.Debug("Add middleware -- martini-contrib/sessions......")
	if "redis" == strings.ToLower(strings.TrimSpace(config.SessionStore)) {
		var err error
		log.Debug("Connect to redis......" + config.Redis.Address)
		if app.Store == nil {
			app.Store, err = redistore.NewRediStoreWithDB(
				config.Redis.Size,
				config.Redis.Network,
				config.Redis.Address,
				config.Redis.Password,
				config.Redis.DB,
				[]byte(config.SessionMask),
			)
			if err != nil {
				log.Error("Failed to connect to redis : " + err.Error())
				return nil
			}
		}
		app.Use(sessions.Sessions(config.SessionName, app.Store))
	} else {
		app.Use(sessions.Sessions(config.SessionName, sessions.NewCookieStore([]byte(config.SessionMask))))
	}

	log.Debug("Connect to databse......")
	if len(config.Database.Database) > 0 {
		app.DbEngine = InitDB(
			config.DBType,
			config.Database.Host,
			config.Database.Port,
			config.Database.Database,
			config.Database.User,
			config.Database.Password,
			config.Database.Verbose,
			config.Database.LogFile,
		)
		if app.DbEngine == nil {
			log.Error("Failed to connect to database (" + config.Database.Database + ")")
			return nil
		}

		// open SQL log or not
		if config.Database.Verbose {
			if len(config.Database.LogFile) > 0 {
				sql_log, err := os.OpenFile(config.Database.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Error("error opening file: %v", err)
					return nil
				}
				app.DbEngine.TraceOn("[SQL]", stdlog.New(sql_log, "sql", stdlog.Lmicroseconds))
			} else {
				app.DbEngine.TraceOn("[SQL]", log.GetLogger("file"))
			}
		} else {
			app.DbEngine.TraceOff()
		}
	}

	log.Debug("Add middleware -- martini-contrib/sessionauth......")
	app.Use(sessionauth.SessionUser(newUser))
	sessionauth.RedirectUrl = config.url["RedirectUrl"]
	sessionauth.RedirectParam = config.url["RedirectParam"]

	log.Debug("Add User defined router mapping......")
	if !cfg.RoterMap(app) {
		log.Error("Failed to add Roter Mapping")
		return nil
	}
	return app
}

func (app *PolarisApplication) Close() bool {
	log.Debug("Uninitializing application......")
	config := app.Config.GetBasicConfig()
	// close the redis session handle
	if app.Store != nil {
		app.Store.Close()
		app.Store = nil
	}
	// close the handle of log file
	if config.LogHandle != nil {
		config.LogHandle.Close()
		config.LogHandle = nil
	}
	return true
}

func (app *PolarisApplication) RotateLog() bool {
	config := app.Config.GetBasicConfig()

	// get new log
	year, month, day := time.Now().Date()
	new_log := config.DirLog + "/" + fmt.Sprintf("app_%04d%02d%02d.log", year, month, day)
	new_handle, err := os.OpenFile(new_log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error("error opening file: %v", err)
		return false
	}

	// replace it
	if config.LogHandle != nil {
		config.LogHandle.Close()
	}
	config.LogHandle = new_handle
	config.LogFile = new_log
	return true
}

func (app *PolarisApplication) RunApp() bool {
	config := app.Config.GetBasicConfig()

	host := os.Getenv("HOST")
	port := config.Port
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
