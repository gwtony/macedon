package goblin

import (
	//"io"
	//"strings"
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/utils"
	"git.lianjia.com/lianjia-sysop/napi/hserver"
)

type AddHandler struct {
	h   *Handler
	mc  *MysqlContext
	log log.Log
}
type DeleteHandler struct {
	h   *Handler
	mc  *MysqlContext
	log log.Log
}
type ReadHandler struct {
	mc  *MysqlContext
	log log.Log
}
type ReadServerHandler struct {
	mc  *MysqlContext
	log log.Log
}
type AddServerHandler struct {
	mc  *MysqlContext
	log log.Log
}
type DeleteServerHandler struct {
	mc  *MysqlContext
	log log.Log
}

func (handler *AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		handler.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Create record request: ", data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" || data.Action == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}
	if data.Expire <= 0 {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}
	/* only support forbidden */
	data.Action = "forbidden"

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	servers, err := handler.mc.QueryGetServer(db, data.Product)
	if err != nil { //TODO: error should be 400 product not found
		hserver.ReturnError(w, err, handler.log)
		return
	}

	data.Ruleid, err = utils.NewUUID()
	if err != nil { //TODO: error should be 500 generate ruleid faild
		hserver.ReturnError(w, err, handler.log)
		return
	}

	data.Expire = int(time.Now().Unix() + int64(data.Expire))
	err = handler.mc.QueryInsert(db, data.Ip, data.Uid, data.Uuid, data.Product, data.Expire, data.Action, data.Ruleid)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	errFlag := 0
	for _, server := range *servers {
		err = handler.h.RuleAdd(server, data.Ip, data.Uid, data.Uuid, data.Expire, data.Action)
		if err != nil {
			handler.log.Error("Rule add to %s failed, data is %s", server, data)
			errFlag = 1
		}
	}

	if errFlag == 1 {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	err = handler.mc.QueryUpdateResult(db, data.Ruleid)
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

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		handler.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Delete record request: ", data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Ip, Uid, Uuid invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Product invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	servers, err := handler.mc.QueryGetServer(db, data.Product)
	if err != nil { //TODO: error should be 400 product not found
		handler.log.Error("Get server failed: ", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}
	handler.log.Debug("Get server result is: ", servers)

	resp, err := handler.mc.QueryRead(db, data.Ip, data.Uid, data.Uuid, data.Product)
	if err != nil {
		handler.log.Error("Read rule failed: ", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}
	ruleid := resp.Ruleid
	action := resp.Action

	errFlag := 0
	for _, server := range *servers {
		err = handler.h.RuleDelete(server, data.Ip, data.Uid, data.Uuid, action)
		if err != nil {
			handler.log.Error("Rule delete to %s failed, data is %s", server, data)
			errFlag = 1
		}
	}

	if errFlag == 1 {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	err = handler.mc.QueryUpdateDeleted(db, ruleid)
	if err != nil {
		handler.log.Error("Update deleted failed: ", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}

func (handler *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

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

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request from %s ", req.RemoteAddr, data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Ip, Uid, Uuid invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Product invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}
	//TODO: Deal wildcard */
	//name := data.Name
	//if strings.Contains(data.Name, "*") {
	//	name = strings.Replace(data.Name, "*", "%", -1)
	//	like = 1
	//}
	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	cresps, err := handler.mc.QueryRead(db, data.Ip, data.Uid, data.Uuid, data.Product)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	//resps := &GoblinReadResponse{}
	//for _, cresp := range *cresps {
	//	resp := GoblinRequest{}
	//	resp.Ip      = cresp.Ip
	//	resp.Uid     = cresp.Uid
	//	resp.Uuid    = cresp.Uuid
	//	resp.Product = cresp.Product
	//	resp.Expire  = cresp.Expire
	//	resp.Action  = resp.Action
	//	*resps = append(*resps, resp)
	//}

	hserver.ReturnResponse(w, cresps, handler.log)
}

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
