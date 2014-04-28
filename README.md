Polaris
=========

Polaris is my application starter in golang. It's not a framework but a integrated solution to mix in several existing components online. 

#### Features ####
-  As you can see, Polaris is highly depends on [go-martini/martini](https://github.com/go-martini/martini) and its official extentions[martini-contrib](https://github.com/martini-contrib/);
-  For data access, Polaris leverage [robinmin/gorp](https://github.com/robinmin/gorp) which is an enhanced version based on [coopernurse/gorp](https://github.com/coopernurse/gorp). It's an ORM-like databse access driver to enable the end user can access their database easily; it compact with the build in "database/sql", so it's easy to support your databse in case it's not supported so far;
-  Meanwhile, Polaris also embedded with the Redis driver so that you can store your session information and cache immediately;

#### Used Packages ####
Polaris try to integrated the following existing projects as the major compoents. Normally, you should set the environment variable GOPATH before you execute the following commands:
```bash
go get -u github.com/go-martini/martini
go get -u github.com/martini-contrib/binding
go get -u github.com/martini-contrib/encoder
go get -u github.com/martini-contrib/render
go get -u github.com/martini-contrib/throttle
go get -u github.com/martini-contrib/strict
go get -u github.com/martini-contrib/secure
go get -u github.com/martini-contrib/csrf
go get -u github.com/martini-contrib/accessflags
go get -u github.com/martini-contrib/gzip
go get -u github.com/martini-contrib/sessionauth
go get -u github.com/martini-contrib/cors
go get -u github.com/martini-contrib/oauth2
go get -u github.com/martini-contrib/acceptlang
go get -u github.com/martini-contrib/logstasher
go get -u github.com/martini-contrib/auth
go get -u github.com/martini-contrib/sessions
go get -u github.com/martini-contrib/method
go get -u github.com/martini-contrib/strip

go get -u github.com/robinmin/gorp
go get -u github.com/mattn/go-adodb
go get -u github.com/mattn/go-sqlite3
go get -u github.com/boj/redistore

go get -u github.com/robinmin/logo

```

