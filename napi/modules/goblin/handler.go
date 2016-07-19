package goblin

import (
	"fmt"
	"bytes"
	//"time"
	"net/http"
	"git.lianjia.com/lianjia-sysop/napi/variable"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type Handler struct {
	add_loc    string
	delete_loc string
	log        log.Log
}

func InitHandler(loc string, log log.Log) *Handler {
	h := &Handler{}
	h.log = log
	h.add_loc = loc + GOBLIN_ADD_LOCATION
	h.delete_loc = loc + GOBLIN_DELETE_LOCATION

	return h
}

func (h *Handler) RuleOperate(addr string, args *bytes.Buffer, op int) error {
	var err error
	var resp *http.Response

	switch op {
	case ADD_RULE:
		h.log.Debug("add rule args is ", args)
		resp, err = http.Post("http://" + addr + h.add_loc, variable.DEFAULT_CONTENT_HEADER, args)
		//resp, err = http.Get("http://" + addr + h.add_loc + "?" + args)
		break
	case DELETE_RULE:
		h.log.Debug("delete rule args is ", args)
		resp, err = http.Post("http://" + addr + h.delete_loc, variable.DEFAULT_CONTENT_HEADER, args)
		//resp, err = http.Get("http://" + addr + h.delete_loc + "?" + args)
		break
	default: /* Should not reach here */
		h.log.Error("Unknown operate code: ", op)
		return errors.InternalServerError
	}

	if err != nil {
		h.log.Error("Opereate service to nginx failed: ", err)
		return errors.BadGatewayError
	}

	defer resp.Body.Close()

	if resp.StatusCode != variable.HTTP_OK {
		h.log.Error("Opereate http status error: %d", resp.StatusCode)
		return errors.BadGatewayError
	}

	return nil
}

func (h *Handler) RuleAdd(addr, ip, uid, uuid string, expire int, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	//if uid != "" {
	//	args = fmt.Sprint(args, "&uid=", uid)
	//}
	if uuid != "" {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}
	args = fmt.Sprint(args, "&expire=", expire, "&punish=", action, "&punish_arg=0\r\n")
	h.log.Error(args)

	/* post data */
	data := bytes.NewBufferString(args)

	return h.RuleOperate(addr, data, ADD_RULE)
	//return h.RuleOperate(addr, args, ADD_RULE)
}

func (h *Handler) RuleDelete(addr, ip, uid, uuid, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	//if uid != "" {
	//	args = fmt.Sprint(args, "&uid=", uid)
	//}
	if uuid != "" {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}

	args = fmt.Sprint(args, "&punish=", action)

	h.log.Error(args)
	data := bytes.NewBufferString(args)
	return h.RuleOperate(addr, data, DELETE_RULE)
	//return h.RuleOperate(addr, args, DELETE_RULE)
}
