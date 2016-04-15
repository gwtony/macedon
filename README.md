# macedon-go
Description
===========
Http api for consul
* create a record
* delete a record
* read a record



Config Sample
=============

```
[default]
addr: host:ip

log: macedon.log
level: debug

location: /dns
purge_ips: "192.168.0.1"
purge_cmd: "purge dns"
ssh_key: /username/.ssh/id_rsa
ssh_port: 22
ssh_user: username
ssh_timeout: 20

consul_addrs: consul_server
domain: domain
```

Schema
=====
* [schema infomation](SCHEMA.md)

Dependency
==========

* [log4go](http://code.google.com/p/log4go)
* [goconfig](https://github.com/msbranco/goconfig)
* [golang/x/ssh](http://golang.org/x/crypto/ssh)
* [mattn/go-getopt](http://github.com/mattn/go-getopt)
