package macedon

import (
	"io"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type CreateHandler struct {
	hs  *HttpServer
	log *Log
}
type DeleteHandler struct {
	hs  *HttpServer
	log *Log
}
type UpdateHandler struct {
	hs  *HttpServer
	log *Log
}
type ReadHandler struct {
	hs  *HttpServer
	log *Log
}
type NotifyHandler struct {
	hs  *HttpServer
	log *Log
}

func returnError(w http.ResponseWriter, resp *Response, err error, log *Log) {
	if err == NoContentError {
		log.Debug("Request no content")
		http.Error(w, "", http.StatusNoContent)
		return
	}
	if err == BadRequestError {
		log.Debug("Return bad request")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err == BadGatewayError {
		log.Debug("Return bad gateway")
		http.Error(w, "", http.StatusBadGateway)
		return
	}

	log.Debug("Return internal server error")
	http.Error(w, "", http.StatusInternalServerError)
}

func returnResponse(w http.ResponseWriter, resp *Response, log *Log) {
	respj, err := json.Marshal(resp)
	if err != nil {
		log.Error("Encode json failed: ", resp)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	log.Debug("Return OK")

	io.WriteString(w, string(respj))
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &Request{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Create record request: ", data)

	/* Check input */
	if data.Name == "" || data.Type == "" || data.Zname == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name, Zname or Type invalid", http.StatusBadRequest)
		return
	}
	if data.Domain_id < 0 || data.Ttl <= 0 || len(data.Records) <= 0 {
		h.log.Error("Domain_id, ttl or records invalid")
		http.Error(w, "Domain_id, ttl, records maybe invalid", http.StatusBadRequest)
		return
	}

	if !strings.EqualFold(data.Type, "a") &&
		!strings.EqualFold(data.Type, "cname") &&
		!strings.EqualFold(data.Type, "ptr") {
		h.log.Error("Type invalid: %s", data.Type)
		http.Error(w, "Type invalid", http.StatusBadRequest)
		return
	}

	if data.Records[0].Content == "" {
		h.log.Error("Empty content in records")
		http.Error(w, "Records invalid", http.StatusBadRequest)
		return
	}

	mc := h.hs.s.MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer mc.Close(db)

	rec := data.Records[0]
	resp, err := mc.QueryCreate(db, data.Name, data.Type, rec.Content, data.Domain_id, data.Ttl)

	if err != nil {
		returnError(w, resp, err, h.log)
		return
	}

	go h.hs.s.dns.UpdateCreate(data.Zname, data.Name, data.Type, rec.Content, data.Ttl)

	returnResponse(w, resp, h.log)
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var resp *Response

	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &Request{}
	json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Delete record request: ", data)

	/* Check input */
	if data.Name == "" || data.Type == "" || data.Zname == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name, Zname or Type invalid", http.StatusBadRequest)
		return
	}
	if !strings.EqualFold(data.Type, "a") &&
		!strings.EqualFold(data.Type, "cname") &&
		!strings.EqualFold(data.Type, "ptr") {
		h.log.Error("Type invalid: %s", data.Type)
		http.Error(w, "Type invalid", http.StatusBadRequest)
		return
	}
	if len(data.Records) > 0 && data.Records[0].Content == "" {
		h.log.Error("Empty content in records")
		http.Error(w, "Records invalid", http.StatusBadRequest)
		return
	}

	mc := h.hs.s.MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer mc.Close(db)

	content := ""
	if len(data.Records) > 0 {
		content = data.Records[0].Content
	}

	resp, err = mc.QueryDelete(db, data.Name, data.Type, content)

	if err != nil {
		returnError(w, resp, err, h.log)
		return
	}

	go func() {
		if h.hs.s.dns != nil {
			h.hs.s.dns.UpdateDelete(data.Zname, data.Name, data.Type, content)
		}
		if h.hs.s.pc != nil {
			go h.hs.s.pc.DoPurge(h.hs.s.sc)
		}
	}()

	returnResponse(w, resp, h.log)
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &Request{}
	json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Update record request: ", data)

	/* Check input */
	if data.Name == "" || data.Type == "" || data.Zname == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name, Zname or Type invalid", http.StatusBadRequest)
		return
	}
	if data.Ttl <= 0 {
		h.log.Error("Rtl invalid")
		http.Error(w, "Rtl invalid", http.StatusBadRequest)
		return
	}
	if !strings.EqualFold(data.Type, "a") &&
		!strings.EqualFold(data.Type, "cname") &&
		!strings.EqualFold(data.Type, "ptr") {
		h.log.Error("Type invalid: %s", data.Type)
		http.Error(w, "Type invalid", http.StatusBadRequest)
		return
	}
	if len(data.Records) == 0 || data.Records[0].Content == "" {
		h.log.Error("Empty content in records")
		http.Error(w, "Records invalid", http.StatusBadRequest)
		return
	}
	if data.Records[0].Disabled != 0 && data.Records[0].Disabled != 1 {
		h.log.Error("Record disabled state invalid")
		http.Error(w, "Record disabled state invalid", http.StatusBadRequest)
		return
	}

	mc := h.hs.s.MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer mc.Close(db)

	rec := data.Records[0]
	resp, err := mc.QueryUpdate(db, data.Name, data.Type, rec.Content, rec.Disabled)
	if err != nil {
		returnError(w, resp, err, h.log)
		return
	}

	go func() {
		if h.hs.s.dns != nil {
			if (rec.Disabled == 1) {
				h.hs.s.dns.UpdateDelete(data.Zname, data.Name, data.Type, rec.Content)
			} else {
				h.hs.s.dns.UpdateCreate(data.Zname, data.Name, data.Type, rec.Content, data.Ttl)
			}
		}
		if h.hs.s.pc != nil {
			go h.hs.s.pc.DoPurge(h.hs.s.sc)
		}
	}()

	returnResponse(w, resp, h.log)
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	like := 0
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &Request{}
	json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s ", req.RemoteAddr, data)

	/* Check input */
	if data.Name == "" || data.Type == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name or Type invalid", http.StatusBadRequest)
		return
	}
	if !strings.EqualFold(data.Type, "a") &&
		!strings.EqualFold(data.Type, "cname") &&
		!strings.EqualFold(data.Type, "ptr") {
		h.log.Error("Type invalid: %s", data.Type)
		http.Error(w, "Type invalid", http.StatusBadRequest)
		return
	}

	/* Deal wildcard */
	name := data.Name
	if strings.Contains(data.Name, "*") {
		name = strings.Replace(data.Name, "*", "%", -1)
		like = 1
	}

	mc := h.hs.s.MysqlContext()
	db, err := mc.Open()
	if err != nil {
		h.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer mc.Close(db)

	resp, err := mc.QueryRead(db, name, data.Type, like)
	resp.Result.Data.Name = data.Name

	if err != nil {
		returnError(w, resp, err, h.log)
		return
	}

	returnResponse(w, resp, h.log)
}

func (h *NotifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &Request{}
	json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Notify record request: ", data)

	/* Check input */
	if data.Zname == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Zname invalid", http.StatusBadRequest)
		return
	}

	mc := h.hs.s.MysqlContext()
	db, err := mc.Open()
	if err != nil {
		mc.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer mc.Close(db)

	resp, err := mc.QueryNotify(db, data.Zname)

	if err != nil {
		returnError(w, resp, err, h.log)
		return
	}

	returnResponse(w, resp, h.log)
}
