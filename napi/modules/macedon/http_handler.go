package macedon

import (
	//"io"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/utils"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type AddHandler struct {
	h      *Handler
	domain string
	log    log.Log
}
type DeleteHandler struct {
	h      *Handler
	domain string
	log    log.Log
}
type ReadHandler struct {
	h      *Handler
	domain string
	log    log.Log
}
type ReadServerHandler struct {
	log log.Log
}
type AddServerHandler struct {
	log log.Log
}
type DeleteServerHandler struct {
	log log.Log
}

func (handler *AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Create record request: ", data)

	/* Check input */
	if data.Name == "" || data.Addr == "" {
		handler.log.Error("Name or addr invalid")
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

    if !strings.HasSuffix(data.Name, handler.domain) {
        handler.log.Error("Domain in name invalid")
        http.Error(w, "Domain in name invalid", http.StatusBadRequest)
        return
    }

	/* "name1.domain.com" to "com/domain/name1" */
	arr := strings.Split(data.Name, ".")
	total := len(arr)
	for i := 0; i < total / 2; i++ {
		tmp := arr[i]
		arr[i] = arr[total - 1 - i]
		arr[total - 1 - i] = tmp
	}
	rec := strings.Join(arr, "/")

	rec = rec + "/" + fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))
	handler.log.Debug("Rec is %s", rec)

	_, err = handler.h.Add(rec, data.Addr, data.Ttl)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}

func (handler *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Create record request: ", data)

	/* Check input */
	if data.Name == "" || data.Addr == "" {
		handler.log.Error("Name or addr invalid")
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

    if !strings.HasSuffix(data.Name, handler.domain) {
        handler.log.Error("Domain in name invalid")
        http.Error(w, "Domain in name invalid", http.StatusBadRequest)
        return
    }

	/* "name1.domain.com" to "com/domain/name1" */
	arr := strings.Split(data.Name, ".")
	total := len(arr)
	for i := 0; i < total / 2; i++ {
		tmp := arr[i]
		arr[i] = arr[total - 1 - i]
		arr[total - 1 - i] = tmp
	}
	rec := strings.Join(arr, "/")

	//rec = rec + "/" + fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))
	handler.log.Debug("Rec is %s", rec)

	resp, err := handler.h.Read(rec)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	found := ""
	if len(resp.Node.Nodes) > 0 {
        resp_rec := &RecValue{}
        for _, v := range resp.Node.Nodes {
			found = v.Key
            json.Unmarshal([]byte(v.Value), &resp_rec)
			if strings.Compare(resp_rec, data.Addr) == 0 {
				break
			}
        }
    }

	if found == "" {
		hserver.ReturnError(w, errors.NoContent, handler.log)
		return
	}

	_, err = handler.h.Delete(found)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}
	//TODO: do purge
	hserver.ReturnResponse(w, nil, handler.log)
}

func (handler *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &MacedonRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Create record request: ", data)

	/* Check input */
	if data.Name == "" || data.Addr == "" {
		handler.log.Error("Name or addr invalid")
	}
	if data.Ttl <= 0 {
		data.Ttl = DEFAULT_TTL
	}

    if !strings.HasSuffix(data.Name, handler.domain) {
        handler.log.Error("Domain in name invalid")
        http.Error(w, "Domain in name invalid", http.StatusBadRequest)
        return
    }

	/* "name1.domain.com" to "com/domain/name1" */
	arr := strings.Split(data.Name, ".")
	total := len(arr)
	for i := 0; i < total / 2; i++ {
		tmp := arr[i]
		arr[i] = arr[total - 1 - i]
		arr[total - 1 - i] = tmp
	}
	rec := strings.Join(arr, "/")

	handler.log.Debug("Rec is %s", rec)

	resp, err := handler.h.Read(rec)

	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	eresp := &MacedonResponse{}
    eresp.Name = data.Name
	if len(resp.Node.Nodes) > 0 {
        resp_rec := &RecValue{}
        for _, v := range resp.Node.Nodes {
            json.Unmarshal([]byte(v.Value), &resp_rec)
			addr := MacedonAddress{}
			addr.Addr = resp_rec.Addr
			addr.Ttl = resp_rec.Ttl
			eresp.Addrs = append(eresp.Addrs, addr)
        }
    }

	hserver.ReturnResponse(w, eresp, handler.log)
}

/*
func (handler *ReadServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request from %s, data is %s", req.RemoteAddr, data)
	if data.Addr == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Addr invalid", http.StatusBadRequest)
		return
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	resps, err := handler.mc.QueryReadServer(db, data.Addr)
	if err != nil {
		handler.log.Error("Query read server failed: %s", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, resps, handler.log)
}

func (handler *AddServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed:%s ", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request from %s ", req.RemoteAddr, data)
	if data.Addr == "" || data.Product == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Addr or product invalid", http.StatusBadRequest)
		return
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	err = handler.mc.QueryAddServer(db, data.Addr, data.Product)
	if err != nil {
		handler.log.Error("Query add server failed: %s", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}

func (handler *DeleteServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request from %s ", req.RemoteAddr, data)
	if data.Addr == "" || data.Product == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Addr or product invalid", http.StatusBadRequest)
		return
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	err = handler.mc.QueryDeleteServer(db, data.Addr, data.Product)
	if err != nil {
		handler.log.Error("Query delete server failed: %s", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}
*/
