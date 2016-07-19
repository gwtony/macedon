package macedon

import (
	"fmt"
	//"bytes"
	//"time"
	"net/http"
	"git.lianjia.com/lianjia-sysop/napi/variable"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type Handler struct {
	eaddr      []string
	eaddr_size int
	loc        string
	log        log.Log
}

func InitHandler(eaddr []string, loc string, log log.Log) *Handler {
	h := &Handler{}
	h.eaddr = eaddr
	h.eaddr_size = len(eaddr)
	h.log = log
	h.loc = loc

	return h
}

func (h *Handler) Operate(loc string, args string, op int) (*EtcdResponse, error) {
	var err error
	var resp *http.Response

	retry := 0

next:
	switch op {
	case ADD:
		h.log.Debug("add record args is %s", args)
		val := &url.Values{}
		val.Add("value", args)
		data := bytes.NewBufferString(val.Encode())
		req, _ := http.NewRequest("PUT", "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc , data)
		req.Header.Set(CONTENT_HEADER, ETCD_CONTENT_HEADER)
		client := &http.Client{}
		resp, err := client.Do(req)
		break
	case DELETE:
		h.log.Debug("delete record args is %s", loc)
		req, _ := http.NewRequest("DELETE", "http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc)
		req.Header.Set(CONTENT_HEADER, ETCD_CONTENT_HEADER)
		client := &http.Client{}
		resp, err := client.Do(req)
		break
	case READ:
		h.log.Debug("Read record loc is %s", loc)
		resp, err = http.Get("http://" + h.eaddr[retry] + DEFAULT_SKYDNS_LOC + loc)
	default: /* Should not reach here */
		h.log.Error("Unknown operate code: ", op)
		return nil, errors.InternalServerError
	}

	if err != nil {
		h.log.Error("Opereate service to etcd failed: ", err)
		retry++
		if retry >= h.eaddr_size {
			return nil, errors.BadGatewayError
		}
		goto next
	}

	if resp.StatusCode != variable.HTTP_OK {
		h.log.Error("Opereate http status error: %d", resp.StatusCode)
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
		cc.log.Error("Read operate result body failed: ", err)
		return nil, InternalServerError
	}
	defer resp.Body.Close()

	eresp := &EtcdResponse{}
	err = json.Unmarshal(body, &eresp)
	if err != nil {
		h.log.Error("Unmarshl from response failed")
		return nil, errors.InternalServerError
	}

	return eresp, nil
}

func (h *Handler) Add(rec, addr string, ttl int) error {
	r := &RecValue{}
	r.Host = addr
	r.Ttl = ttl
	b, err := json.Marshal(r)
	if err != nil {
		h.log.Error("Add marshal failed")
		return err
	}
	args := string(b)

	h.log.Info(args)

	return h.Operate(rec, args, ADD)
}

func (h *Handler) Delete(rec) error {
	return h.Operate(rec, "", DELETE)
}

func (h *Handler) Read(rec) error {
	return h.Operate(rec, "", READ)
}
