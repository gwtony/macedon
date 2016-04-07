# macedon-go
Description
===========
Mysql api for powerdns
* create a record
* delete a record
* update a record
* read a record
* notify zone update



Config Sample
=============

```
[default]
addr: host:ip
mysql_addr: mysql_host:3306
mysql_dbname: dbname
mysql_dbuser: dbuser
mysql_dbpwd: dbpwd

log: macedon.log
level: debug

location: /dns
purge_ips: "192.168.0.1"
purge_cmd: "purge dns"
ssh_key: /username/.ssh/id_rsa
ssh_port: 22
ssh_user: username
ssh_timeout: 20
```

Schema
=====
* [schema infomation](SCHEMA.md)

Dependency
==========

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [golang/x/ssh](http://golang.org/x/crypto/ssh)
* [go-sql-driver/mysql](http://github.com/go-sql-driver/mysql)
