# macedon-go
Description
===========
Macedon is a mysql api for powerdns
* create a record
* delete a record
* update a record
* read a record



Config Sample
=============

```
addr: host:port
maddr: mysql_host:port
dbname: database_name
dbuser: database_user
dbpwd: data_password

log: file to log
level: debug

create_location: /create
delete_location: /delete
update_location: /update
read_location: /read
```

Dependency
==========

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [golang/x/ssh](http://golang.org/x/crypto/ssh)
