package macedon

import (
	//"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/gwtony/gapi/variable"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/errors"
)

// Handler Etcd handler
type Handler struct {
	eaddr      []string
	eaddrSize  int
	loc        string
	log        log.Log
}

// InitHandler inits handler
func InitHandler(eaddr []string, loc string, log log.Log) *Handler {
	h := &Handler{}
	h.eaddr = eaddr
	h.eaddrSize = len(eaddr)
	h.log = log
	h.loc = loc

	return h
}

// Operate operates etcd
func (h *Handler) Operate(loc, args string, op int, arpa, server, scan bool) (*EtcdResponse, error) {
	var err error
	var floc string
	var resp *http.Response

	retry := 0

next:
	switch op {
	case ADD:
		h.log.Debug("add record args is %s", args)
		val := &url.Values{}
		val.Add("value", args)
		data := bytes.NewBufferString(val.Encode())
		if arpa {
			floc = "http://" + h.eaddr[retry] + DEFAULT_ARPA_LOC + loc
		} else {
			floc = "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc
		}
		if server {
			floc = "http://" + h.eaddr[retry] + DEFAULT_PURGE_SERVER_LOC + loc
		}
		h.log.Debug("add url is %s", floc)
		req, _ := http.NewRequest("PUT", floc , data)
		req.Header.Set(CONTENT_HEADER, ETCD_CONTENT_HEADER)
		client := &http.Client{}
		resp, err = client.Do(req)
		break
	case DELETE:
		if arpa {
			floc = "http://" + h.eaddr[retry] + DEFAULT_ARPA_LOC + loc
		} else {
			floc = "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc
		}
		if server {
			floc = "http://" + h.eaddr[retry] + DEFAULT_PURGE_SERVER_LOC + loc
		}
		h.log.Debug("del url is %s", floc)
		req, _ := http.NewRequest("DELETE", floc, nil)
		req.Header.Set(CONTENT_HEADER, ETCD_CONTENT_HEADER)
		client := &http.Client{}
		resp, err = client.Do(req)
		break
	case READ:
		if arpa {
			floc = "http://" + h.eaddr[retry] + DEFAULT_ARPA_LOC + loc
		} else {
			if scan {
				floc = "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc + DEFAULT_SCAN_ARGS
			} else {
				floc = "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc
			}
		}
		if server {
			floc = "http://" + h.eaddr[retry] + DEFAULT_PURGE_SERVER_LOC
		}
		h.log.Debug("read url is %s", floc)
		resp, err = http.Get(floc)
		break
	default: /* Should not reach here */
		h.log.Error("Unknown operate code: ", op)
		return nil, errors.InternalServerError
	}

	if err != nil {
		h.log.Error("Opereate service to etcd failed: ", err)
		retry++
		if retry >= h.eaddrSize {
			return nil, errors.BadGatewayError
		}
		goto next
	}

	defer resp.Body.Close()

	if resp.StatusCode != variable.HTTP_OK && resp.StatusCode != variable.HTTP_CREATED {
		h.log.Debug("Opereate http status error: %d", resp.StatusCode)
		if (op != ADD) && (resp.StatusCode == variable.HTTP_NOT_FOUND) {
			return nil, errors.NoContentError
		}
		return nil, errors.BadGatewayError
	}

	if op == ADD || op == DELETE {
		return nil, nil
	}

	/* read operation */
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.log.Error("Read operate result body failed: ", err)
		return nil, errors.InternalServerError
	}
	//defer resp.Body.Close()

	eresp := &EtcdResponse{}
	err = json.Unmarshal(body, &eresp)
	if err != nil {
		h.log.Error("Unmarshl from response failed")
		return nil, errors.InternalServerError
	}

	return eresp, nil
}

// Add adds record to etcd
func (h *Handler) Add(rec, addr string, ttl int, arpa, server bool) (*EtcdResponse, error) {
	r := &RecValue{}
	r.Host = addr
	r.Ttl = ttl
	b, err := json.Marshal(r)
	if err != nil {
		h.log.Error("Add marshal failed")
		return nil, err
	}
	args := string(b)

	h.log.Info(args)

	return h.Operate(rec, args, ADD, arpa, server, false)
}

// Delete deletes record to etcd
func (h *Handler) Delete(rec string, arpa, server bool) (*EtcdResponse, error) {
	return h.Operate(rec, "", DELETE, arpa, server, false)
}

// Read reads record from etcd
func (h *Handler) Read(rec string, arpa, server, scan bool) (*EtcdResponse, error) {
	return h.Operate(rec, "", READ, arpa, server, scan)
}
