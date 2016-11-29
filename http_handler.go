package macedon

import (
	"net"
	"fmt"
	"time"
	"math/rand"
	"strings"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

// AddHandler Add record handler
type AddHandler struct {
	h      *Handler
	domain string
	token  string
	pc     *PurgeContext
	log    log.Log
}

// DeleteHandler Delete record handler
type DeleteHandler struct {
	h      *Handler
	domain string
	token  string
	pc     *PurgeContext
	log    log.Log
}

// ReadHandler Read record handler
type ReadHandler struct {
	h      *Handler
	domain string
	token  string
	log    log.Log
}

// ScanHandler Scan record handler
type ScanHandler struct {
	h      *Handler
	domain string
	token  string
	log    log.Log
}

// ReadServerHandler Read server record handler
type ReadServerHandler struct {
	h     *Handler
	token string
	log   log.Log
}

// AddServerHandler Add server record handler
type AddServerHandler struct {
	h     *Handler
	pc    *PurgeContext
	token string
	log   log.Log
}

// DeleteServerHandler Delete server record handler
type DeleteServerHandler struct {
	h     *Handler
	pc    *PurgeContext
	token string
	log   log.Log
}

// splitReverse Split and reverse "a.b.c.d" to ['d', 'c', 'b', 'a']
func splitReverse(str, sep string) []string {
	arr := strings.Split(str, sep)
	total := len(arr)
	for i := 0; i < total / 2; i++ {
		tmp := arr[i]
		arr[i] = arr[total - 1 - i]
		arr[total - 1 - i] = tmp
	}
	return arr
}

// ServeHTTP router interface
func (handler *AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isArpa bool
	var arr []string
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}
	handler.log.Info("Add record request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}

	/* Check input */
	if data.Name == "" || data.Address == "" {
		api.ReturnError(r, w, errors.Jerror("Name or address invalid"), errors.BadRequestError, log)
		return
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = splitReverse(data.Name, ".")
	}
	rec := strings.Join(arr, "/")

	resp, err := handler.h.Read(rec, isArpa, false, false)
	if err != nil && err != errors.NoContentError {
		api.ReturnError(r, w, errors.Jerror("Read record failed"), err, log)
		return
	}

	// check exists record
	if resp != nil {
		exist := false
		if len(resp.Node.Nodes) > 0 {
			respRec := &RecValue{}
			for _, v := range resp.Node.Nodes {
				json.Unmarshal([]byte(v.Value), &respRec)
				if strings.Compare(respRec.Host, data.Address) == 0 {
					exist = true
					break
				}
			}
		} else {
			if resp.Node.Value != "" {
				respRec := RecValue{}
				json.Unmarshal([]byte(resp.Node.Value), &respRec)
				if strings.Compare(respRec.Host, data.Address) == 0 {
					exist = true
				}
			}
		}

		if exist {
			api.ReturnError(r, w, errors.Jerror("Record exists"), errors.ConflictError, log)
			return
		}
	}
	// one arpa record should only match to one domain name, one to multi not supported by skydns
	if !isArpa {
		rec = rec + "/" + fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))
	}

	_, err = handler.h.Add(rec, data.Address, data.Ttl, isArpa, false)

	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Add record failed"), err, log)
		return
	}

	go handler.pc.DoPurge(data.Name)

	api.ReturnResponse(r, w, "", log)
}

// ServeHTTP router interface
func (handler *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isArpa bool
	var arr []string
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}
	handler.log.Info("Delete record request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}
	/* Check input */
	if data.Name == "" || data.Address == "" {
		api.ReturnError(r, w, errors.Jerror("Name or address invalid"), errors.BadRequestError, log)
		return
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = splitReverse(data.Name, ".")
	}
	rec := strings.Join(arr, "/")

	resp, err := handler.h.Read(rec, isArpa, false, false)
	if err != nil {
		if err == errors.NoContentError {
			api.ReturnError(r, w, errors.Jerror("No record found"), err, log)
		} else {
			api.ReturnError(r, w, errors.Jerror("Read record failed"), err, log)
		}
		return
	}

	var found []string
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
			json.Unmarshal([]byte(v.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				if isArpa {
					found = append(found, strings.TrimPrefix(v.Key, DEFAULT_TRIM_ARPA_KEY))
				} else {
					found = append(found, strings.TrimPrefix(v.Key, DEFAULT_TRIM_KEY))
				}
			}
		}
	} else {
		if resp.Node.Value != "" {
			respRec := RecValue{}
			json.Unmarshal([]byte(resp.Node.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				if isArpa {
					found = append(found, strings.TrimPrefix(resp.Node.Key, DEFAULT_TRIM_ARPA_KEY))
				} else {
					found = append(found, strings.TrimPrefix(resp.Node.Key, DEFAULT_TRIM_KEY))
				}
			}
		}
	}

	if len(found) == 0 {
		api.ReturnError(r, w, errors.Jerror("No record found"), errors.NoContentError, log)
		return
	}

	for _, v := range found {
		_, err = handler.h.Delete(v, isArpa, false)
		if err != nil {
			api.ReturnError(r, w, errors.Jerror("Delete record failed"), err, log)
			return
		}
	}

	go handler.pc.DoPurge(data.Name)

	api.ReturnResponse(r, w, "", log)
}

// parseResponse Parse etcd response to macedon response
func parseResponse(key string, n *Node, mresp *MacedonResponse) {
	if n.Dir {
		for _, v := range n.Nodes {
		parseResponse(n.Key, &v, mresp)
		}
	} else {
		respRec := &RecValue{}
		json.Unmarshal([]byte(n.Value), &respRec)
		r := MacedonResponseRecord{}
		rkey := splitReverse(n.Key, "/")
		r.Name = strings.Join(rkey[1: len(rkey) - 2], ".") // Get real domain name
		r.Address = respRec.Host
		r.Ttl = respRec.Ttl
		mresp.Result = append(mresp.Result, r)
	}
}

// ServeHTTP router interface
func (handler *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isArpa bool
	var arr []string
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}
	handler.log.Info("Read record request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}
	/* Check input */
	if data.Name == "" {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = splitReverse(data.Name, ".")
	}
	rec := strings.Join(arr, "/")

	resp, err := handler.h.Read(rec, isArpa, false, false)
	if err != nil {
		if err == errors.NoContentError {
			api.ReturnError(r, w, errors.Jerror("No record found"), err, log)
		} else {
			api.ReturnError(r, w, errors.Jerror("Read record failed"), err, log)
		}
		return
	}

	mresp := &MacedonResponse{}
	if resp.Node.Dir {
		for _, v := range resp.Node.Nodes {
			parseResponse(resp.Node.Key, &v, mresp)
		}
	} else {
		if resp.Node.Value != "" {
			respRec := RecValue{}
			json.Unmarshal([]byte(resp.Node.Value), &respRec)
			r := MacedonResponseRecord{}
			rkey := splitReverse(resp.Node.Key, "/")
			r.Name = strings.Join(rkey[1: len(rkey) - 2], ".")
			r.Address = respRec.Host
			r.Ttl = respRec.Ttl
			mresp.Result = append(mresp.Result, r)
		}
	}

	if len(mresp.Result) == 0 {
		api.ReturnError(r, w, errors.Jerror("No record found"), errors.NoContentError, log)
		return
	}

	respj, err := json.Marshal(mresp)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Encode json failed"), errors.NoContentError, log)
		return
	}
	api.ReturnResponse(r, w, string(respj), log)
}

// ServeHTTP router interface
func (handler *ScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isArpa bool
	var arr []string
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}
	handler.log.Info("Scan record request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}
	/* Check input */
	if data.Name == "" {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	if !(strings.HasSuffix(data.Name, handler.domain) || net.ParseIP(data.Name) != nil) {
		api.ReturnError(r, w, errors.Jerror("Name invalid"), errors.BadRequestError, log)
		return
	}

	isArpa = false
	if net.ParseIP(data.Name) != nil {
		isArpa = true
	}

	if isArpa {
		/* 10.1.1.2 to 10/1/1/2 */
		arr = strings.Split(data.Name, ".")
	} else {
		/* "name1.domain.com" to "com/domain/name1" */
		arr = splitReverse(data.Name, ".")
	}
	rec := strings.Join(arr, "/")

	resp, err := handler.h.Read(rec, isArpa, false, true)
	if err != nil {
		if err == errors.NoContentError {
			api.ReturnError(r, w, errors.Jerror("No record found"), err, log)
		} else {
			api.ReturnError(r, w, errors.Jerror("Read record failed"), err, log)
		}
		return
	}

	mresp := &MacedonResponse{}
	if resp.Node.Dir {
		for _, v := range resp.Node.Nodes {
			parseResponse(resp.Node.Key, &v, mresp)
		}
	} else {
		if resp.Node.Value != "" {
			respRec := RecValue{}
			json.Unmarshal([]byte(resp.Node.Value), &respRec)
			r := MacedonResponseRecord{}
			rkey := splitReverse(resp.Node.Key, "/")
			r.Name = strings.Join(rkey[1: len(rkey) - 2], ".")
			r.Address = respRec.Host
			r.Ttl = respRec.Ttl
			mresp.Result = append(mresp.Result, r)
		}
	}

	if len(mresp.Result) == 0 {
		api.ReturnError(r, w, errors.Jerror("No record found"), errors.NoContentError, log)
		return
	}

	respj, err := json.Marshal(mresp)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Encode json failed"), err, log)
		return
	}
	api.ReturnResponse(r, w, string(respj), log)
}

// ServeHTTP router interface
func (handler *AddServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}

	handler.log.Info("Add server request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}
	if data.Address == "" {
		api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
		return
	}
	if !strings.Contains(data.Address, ":") && net.ParseIP(data.Address) == nil {
		api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
		return
	}
	if strings.Contains(data.Address, ":") {
		tarr := strings.Split(data.Address, ":")
		if net.ParseIP(tarr[0]) == nil {
			api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
			return
		}
	}
	if !strings.Contains(data.Address, ":") {
		data.Address = data.Address + ":" + DEFAULT_PURGE_PORT
	}

	rec := fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))

	_, err = handler.h.Add(rec, data.Address, 0, false, true)

	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Add server failed"), err, log)
		return
	}
	handler.pc.AddServer(data.Address)

	api.ReturnResponse(r, w, "", log)
}

// ServeHTTP router interface
func (handler *ReadServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}

	handler.log.Info("Read server request from client: %s", r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}

	resp, err := handler.h.Read("", false, true, false)

	if err != nil {
		api.ReturnError(r, w, errors.Jerror("No server found"), err, log)
		return
	}

	eresp := &ServerResponse{}
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
		json.Unmarshal([]byte(v.Value), &respRec)
			addr := ServerResponseRecord{}
			addr.Address = respRec.Host
			eresp.Result = append(eresp.Result, addr)
		}
	}

	if len(eresp.Result) == 0 {
		api.ReturnError(r, w, errors.Jerror("No server found"), errors.NoContentError, log)
		return
	}

	respj, err := json.Marshal(eresp)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Encode json failed"), errors.InternalServerError, log)
		return
	}
	api.ReturnResponse(r, w, string(respj), log)
}

// ServeHTTP router interface
func (handler *DeleteServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := handler.log

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, log)
		return
	}

	result, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, log)
		return
	}
	r.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, log)
		return
	}

	handler.log.Info("Delete server request: (%s) from client: %s", data, r.RemoteAddr)

	if handler.token != "" && strings.Compare(data.Token, handler.token) != 0 {
		api.ReturnError(r, w, errors.Jerror("Token invalid"), errors.ForbiddenError, log)
		return
	}
	if data.Address == "" {
		api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
		return
	}
	if !strings.Contains(data.Address, ":") && net.ParseIP(data.Address) == nil {
		api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
		return
	}
	if strings.Contains(data.Address, ":") {
		tarr := strings.Split(data.Address, ":")
		if net.ParseIP(tarr[0]) == nil {
			api.ReturnError(r, w, errors.Jerror("Address invalid"), errors.BadRequestError, log)
			return
		}
	}
	if !strings.Contains(data.Address, ":") {
		data.Address = data.Address + ":" + DEFAULT_PURGE_PORT
	}

	resp, err := handler.h.Read("", false, true, false)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read server failed"), err, log)
		return
	}

	var found []string
	if len(resp.Node.Nodes) > 0 {
		respRec := &RecValue{}
		for _, v := range resp.Node.Nodes {
			json.Unmarshal([]byte(v.Value), &respRec)
			if strings.Compare(respRec.Host, data.Address) == 0 {
				found = append(found, strings.TrimPrefix(v.Key, DEFAULT_TRIM_SERVER_KEY))
			}
		}
	}

	if len(found) == 0 {
		api.ReturnError(r, w, errors.Jerror("No server found"), errors.NoContentError, log)
		return
	}

	for _, v := range found {
		_, err = handler.h.Delete(v, false, true)
		if err != nil {
			api.ReturnError(r, w, errors.Jerror("Delete server failed"), err, log)
			return
		}
	}

	handler.pc.DeleteServer(data.Address)

	api.ReturnResponse(r, w, "", log)
}
