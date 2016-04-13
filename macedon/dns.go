package macedon

import (
	"fmt"
	"time"
	"strings"
	"github.com/miekg/dns"
)

type DnsUpdater struct {
	servers []string
	slen    int
	timeout time.Duration
	log     *Log
}

func InitDnsUpdater(servers string, log *Log) (*DnsUpdater, error){
	du := &DnsUpdater{}
	du.servers = strings.Split(servers, ",")
	du.slen = len(du.servers)
	if du.slen <= 0 {
		return nil, BadConfigError
	}
	for idx, srv := range du.servers {
		if !strings.Contains(srv, ":") {
			du.servers[idx] = srv + ":53"
		}
	}
	log.Debug(du.servers)
	du.log = log

	return du, nil
}

func (du *DnsUpdater) updateRecord(host, zname, name, types, content string, ttl int, cord bool) error {
	var rrs []dns.RR
	m := &dns.Msg{}
	record_str := ""
	if !strings.HasSuffix(zname, ".") {
		zname = zname + "."
	}
	du.log.Debug("zone name: ", zname)

	m.SetUpdate(zname)

	if cord { /* Create record */
		record_str = fmt.Sprintf("%s. %d %s %s", name, ttl, types, content)
	} else {  /* Delete record */
		record_str = fmt.Sprintf("%s. %s") //TODO: content ?
	}

	du.log.Debug(record_str)

	rr, err := dns.NewRR(record_str)
	if err != nil {
		du.log.Error("New rr failed")
		return err
	}

	rrs = append(rrs, rr)
	du.log.Debug(rrs)
	if cord {
		du.log.Debug("Insert record")
		m.Insert(rrs)
	} else {
		du.log.Debug("Remove record")
		m.Remove(rrs)
	}

	_, err = dns.Exchange(m, host)
	if err != nil {
		du.log.Error("Dns exchange error: ", err)
		return err
	}

	return nil
}

func (du *DnsUpdater) updateAll(zname, name, types, content string, ttl int, cord bool) error {
	du.log.Debug("Update all")

	ch := make(chan int, du.slen)

	for _, host := range du.servers {
		du.log.Debug("Update create ip: %s", host)

		go func(ip string) {
			defer func () { ch <- 1 }()

			if cord {
				du.updateRecord(host, zname, name, types, content, ttl, cord)
			} else {
				du.updateRecord(host, zname, name, types, content, ttl, cord)
			}
		}(host)
	}

	for i := 0; i < du.slen; i++ {
		<-ch
	}

	du.log.Debug("All update done")

	return nil
}

func (du *DnsUpdater) UpdateCreate(zname, name, types, content string, ttl int) error {
	du.log.Debug("Update create")

	return du.updateAll(zname, name, types, content, ttl, true)
}
func (du *DnsUpdater) UpdateDelete(zname, name, types, content string) error {
	du.log.Debug("Update delete")

	return du.updateAll(zname, name, types, content, 0, false)
}
