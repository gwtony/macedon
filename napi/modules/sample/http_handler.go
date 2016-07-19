package sample

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
type ClearExpireHandler struct {
	token string
	mc    *MysqlContext
	log   log.Log
}

func (handler *AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var bandstr []string
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

	data := &SampleRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Create record request: ", data)

	/* Check input */
	if data.Type != SAMPLE_TYPE {
		handler.log.Error("Post type arguments invalid")
		http.Error(w, "type invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" {
		handler.log.Error("Post product arguments invalid")
		http.Error(w, "Product invalid", http.StatusBadRequest)
		return
	}
	if len(data.Match.Host) == 0 || len(data.Match.Band) == 0 || data.Match.Expire <= 0 {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Host or band or expire invalid", http.StatusBadRequest)
		return
	}
	if data.Match.Expire <= int(time.Now().Unix()) {
		handler.log.Error("Post expire arguments invalid")
		http.Error(w, "Timestamp expired", http.StatusBadRequest)
		return
	}
	if data.Action.Type != SAMPLE_ACTION_TYPE || data.Action.Value == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Action invalid", http.StatusBadRequest)
		return
	}

	data.Action.Value = strings.TrimSpace(data.Action.Value)
	idx := strings.Index(data.Action.Value, ":")
	if idx <= 0 || idx == (len(data.Action.Value) - 1) {
		handler.log.Error("Post action value arguments invalid")
		http.Error(w, "Action value invalid", http.StatusBadRequest)
		return
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
		hserver.ReturnError(w, err, handler.log)
		return
	}

	data.Ruleid, err = utils.NewUUID()
	if err != nil { //TODO: error should be 500 generate ruleid faild
		hserver.ReturnError(w, err, handler.log)
		return
	}

	for _, v := range data.Match.Band {
		bandstr = append(bandstr, strconv.Itoa(v))
	}

	bands := strings.Join(bandstr, ",")
	hosts := strings.Join(data.Match.Host, ",")

	err = handler.mc.QueryInsert(db, data.Type, hosts, bands,
			data.Action.Type, data.Action.Value, data.Ruleid, data.Product, data.Match.Expire)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

	errFlag := 0
	for _, h := range data.Match.Host {
		for _, b := range data.Match.Band {
			for _, server := range *servers {
				err = handler.h.RuleAdd(server, h, b, data.Match.Expire, data.Action.Value)
				if err != nil {
					handler.log.Error("Add rule to %s failed, host: %s, band: %d", server, h, b)
					errFlag = 1
				}
			}
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

	resp := &SampleAddResponse{}
	resp.Ruleid = data.Ruleid

	hserver.ReturnResponse(w, resp, handler.log)
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

	data := &SampleRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Delete record request: ", data)

	if data.Ruleid == "" {
		handler.log.Error("Post ruleid arguments invalid")
		http.Error(w, "Ruleid invalid", http.StatusBadRequest)
		return
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	resp, err := handler.mc.QueryReadByRuleid(db, data.Ruleid)
	if err != nil {
		handler.log.Error("Read rule failed: ", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	product := resp.Product

	servers, err := handler.mc.QueryGetServer(db, product)
	if err != nil { //TODO: error should be 400 product not found
		handler.log.Error("Get server failed: ", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}
	handler.log.Debug("Get server result is: ", servers)

	errFlag := 0
	for _, h := range resp.Match.Host {
		for _, b := range resp.Match.Band {
			for _, server := range *servers {
				err = handler.h.RuleDelete(server, h, b, resp.Match.Expire, resp.Action.Value)
				if err != nil {
					handler.log.Error("Delete rule to %s failed, host: %s, band: %d", server, h, b)
					errFlag = 1
				}
			}
		}
	}

	if errFlag == 1 {
		hserver.ReturnError(w, errors.BadGatewayError, handler.log)
		return
	}

	err = handler.mc.QueryUpdateDeleted(db, data.Ruleid)
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

	data := &SampleRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	handler.log.Info("Read record request from %s ", req.RemoteAddr, data)

	/* Check input */
	if data.Host == "" {
		handler.log.Error("Post arguments invalid")
		http.Error(w, "Host invalid", http.StatusBadRequest)
		return
	}
	if data.Type == "" || (data.Type != "equal" && data.Type != "like") {
		handler.log.Error("Post type arguments invalid")
		http.Error(w, "Type invalid", http.StatusBadRequest)
		return
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

	cresps, err := handler.mc.QueryReadByHost(db, data.Host, data.Type)
	if err != nil {
		hserver.ReturnError(w, err, handler.log)
		return
	}

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

func (handler *ClearExpireHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	data := &ClearRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		handler.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}

	if data.Token != handler.token {
		handler.log.Error("Token check failed")
		http.Error(w, "Token check failed", http.StatusForbidden)
		return
	}
	if data.Expire >= int(time.Now().Unix()) {
		handler.log.Error("Expire is too large")
		http.Error(w, "Expire should be a past timestamp", http.StatusBadRequest)
		return
	}

	db, err := handler.mc.Open()
	if err != nil {
		handler.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer handler.mc.Close(db)

	err = handler.mc.QueryClearExpire(db, data.Expire)
	if err != nil {
		handler.log.Error("Query clear expires failed: %s", err)
		hserver.ReturnError(w, err, handler.log)
		return
	}

	hserver.ReturnResponse(w, nil, handler.log)
}
