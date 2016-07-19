# napi
Description
===========
Http api for nginx module
* sample
* [goblin]

Installation
============
* make ('dist' dir will be created)
* Mysql table (table description in sql/sample_router.sql)

Config Sample
=============

```
[default]
addr: host:port       #napi server address

log: napi.log         #optional, default is "../log/napi.log"
level: error          #debug, info, error

[sample]
mysql_addr: host:port #mysql address
mysql_dbname: name    #db name
mysql_dbuser: user    #db user
mysql_dbpwd: pwd      #db pwd
location: /sample     #sample nginx module admin location
api_location: /sample #sample api location
```

Usage
=====
* -f config file
* -h help
* -v version

Schema
=====
* [schema infomation](SCHEMA.md) (To be continued)

Dependency
==========

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [golang/x/ssh](http://golang.org/x/crypto/ssh)
* [mysql](https://github.com/go-sql-driver/mysql)
