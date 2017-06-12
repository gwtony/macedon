# macedon
Description
===========
Http (with purge dns) api for etcd
* create a record
* delete a record
* read a record

Installation
============
make

Config Sample
=============

```
[default]
addr: host:ip

log: ../log/macedon.log
level: debug

[macedon]
etcd_addr: 1.1.1.1:2379,1.1.1.2:2379,1.1.1.3:2379
api_location: /macedon
domain: domain
token: some_token
```

Usage
=====
* -f config file
* -h help
* -v version

API Reference
=====
* [API REFERENCE](REFERENCE.md)

Dependency
==========

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [gapi](http://github.com/gwtony/gapi)
