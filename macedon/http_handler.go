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
type ReadHandler struct {
	hs  *HttpServer
	log *Log
}

func returnError(w http.ResponseWriter, err error, log *Log) {
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
	if resp == nil {
		log.Debug("Return OK")
		w.WriteHeader(http.StatusOK)
		return
	}
	respj, err := json.Marshal(resp)
	if err != nil {
		log.Error("Encode json failed: ", resp)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", DEFAULT_CONTENT_HEADER)
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

	if req.Body == nil {
		h.log.Error("Post body not existed")
		http.Error(w, "Post body not existed", http.StatusBadRequest)
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
	if data.Name == "" || data.Address == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name or Address invalid", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(data.Name, h.hs.s.domain) {
		h.log.Error("Post arguments domain invalid")
		http.Error(w, "Domain in name invalid", http.StatusBadRequest)
		return
	}
	name := strings.TrimSuffix(data.Name, h.hs.s.domain)

	err = h.hs.s.cc.RegisterService(name, data.Address)

	if err != nil {
		returnError(w, err, h.log)
		return
	}

	returnResponse(w, nil, h.log)
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	if req.Body == nil {
		h.log.Error("Post body not existed")
		http.Error(w, "Post body not existed", http.StatusBadRequest)
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
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Delete record request: ", data)

	/* Check input */
	if data.Name == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(data.Name, h.hs.s.domain) {
		h.log.Error("Post arguments domain invalid")
		http.Error(w, "Domain in name invalid", http.StatusBadRequest)
		return
	}
	name := strings.TrimSuffix(data.Name, h.hs.s.domain)

	err = h.hs.s.cc.DeRegisterService(name, data.Address)
	if err != nil {
		returnError(w, err, h.log)
		return
	}

	go func() {
		if h.hs.s.pc != nil {
			go h.hs.s.pc.DoPurge(h.hs.s.sc)
		}
	}()

	returnResponse(w, nil, h.log)
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	if req.Body == nil {
		h.log.Error("Post body not existed")
		http.Error(w, "Post body not existed", http.StatusBadRequest)
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
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s ", req.RemoteAddr, data)

	/* Check input */
	if data.Name == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(data.Name, h.hs.s.domain) {
		h.log.Error("Post arguments domain invalid")
		http.Error(w, "Domain in name invalid", http.StatusBadRequest)
		return
	}
	name := strings.TrimSuffix(data.Name, h.hs.s.domain)

	//TODO: Deal wildcard */
	cresps, err := h.hs.s.cc.ListService(name, data.Address)
	if err != nil {
		returnError(w, err, h.log)
		return
	}
	h.log.Debug(name)
	resps := &Response{}
	for _, cresp := range *cresps {
		resp := Request{}
		resp.Name    = data.Name
		resp.Address = cresp.ServiceAddress
		resp.Port    = cresp.ServicePort
		*resps = append(*resps, resp)
	}

	returnResponse(w, resps, h.log)
}
