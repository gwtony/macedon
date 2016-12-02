package handler

import (
	"net"
	"sync"
	"time"
	"strings"
	"encoding/json"
	"github.com/gwtony/gapi/log"
)

// PurgeContext Purge context
type PurgeContext struct {
	handler  *Handler
	lock     *sync.RWMutex

	servers  []string
	lservers int
	cmd      string

	timeout  time.Duration

	log      log.Log
}

// InitPurgeContext inits purge context
func InitPurgeContext(h *Handler, cmd string, timeout time.Duration, log log.Log) *PurgeContext {
	pc := &PurgeContext{}

	pc.handler = h
	pc.cmd = cmd
	pc.timeout = timeout
	pc.log = log

	pc.lservers = 0

	pc.lock = &sync.RWMutex{}

	return pc
}

// DoPurge do purge, do not return any error
func (pc *PurgeContext) DoPurge(name string) {
	pc.log.Debug("Do purge")

	ch := make(chan int, pc.lservers)

	if pc.lservers == 0 {
		resp, err := pc.handler.Read("", false, true, false)
		if err != nil {
			pc.log.Error("Do purge read server failed", err)
			return
		}
		var arr []string
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
			json.Unmarshal([]byte(v.Value), &respRec)
			arr = append(arr, respRec.Host)
		}
		if len(arr) > 0 {
			pc.log.Debug("Update server: ", arr)
			pc.UpdateServer(arr)
		} else {
			return
		}
	}

	pc.lock.RLock()
	for _, host := range pc.servers {
		pc.log.Debug("Purge ip: %s", host)

		go func(addr string) {
			defer func () { ch <- 1 }()
			conn, err := net.DialTimeout("tcp", addr, pc.timeout)
			if err != nil {
				pc.log.Warn("Connect to %s failed", addr)
				return
			}
			defer conn.Close()
			conn.SetWriteDeadline(time.Now().Add(pc.timeout))
			conn.Write([]byte(pc.cmd + " " + name))

		}(host)
	}
	pc.lock.RUnlock()

	for i := 0; i < pc.lservers; i++ {
		<-ch
	}

	pc.log.Debug("All purge done")

	return
}

// AddServer adds server
func (pc *PurgeContext) AddServer(server string) {
	pc.lock.Lock()
	pc.servers = append(pc.servers, server)
	pc.lservers++
	pc.lock.Unlock()
}

// DeleteServer deletes server
func (pc *PurgeContext) DeleteServer(server string) {
	var i, j int
	var arr []string
	j = 0
	pc.lock.Lock()
	for i = 0; i < pc.lservers; i++ {
		if strings.Compare(pc.servers[i], server) == 0 {
			if i - 1 >= 0 {
				arr = append(arr, pc.servers[j:i]...)
			}
			j += i - j + 1
		}
	}
	arr =  append(arr, pc.servers[j:]...)

	pc.servers = arr
	pc.lservers = len(arr)
	pc.lock.Unlock()
}

// UpdateServer updates servers
func (pc *PurgeContext) UpdateServer(servers []string) {
	pc.lock.Lock()
	pc.servers = servers
	pc.lservers = len(servers)
	pc.lock.Unlock()
}
