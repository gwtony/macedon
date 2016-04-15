package macedon

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//var {
//	DEFAULT_HEADER = "application/json;charset=utf-8"
//	REGISTER_LOC   = "/v1/agent/register"
//	UNREGISTER_LOC = "/v1/agent/unregister"
//	READ_LOC       = "/v1/agent/" //TODO:
//}

//const (
//	REG = iota
//	UNREG
//	READ
//)

type ConsulContext struct {
	addrs     []string
	index     int
	addr_num  int

	reg_loc   string
	dereg_loc string
	read_loc  string

	log       *Log
}

func InitConsulContext(addrs_str, reg_loc, dereg_loc, read_loc string, log *Log) (*ConsulContext, error) {
	cc := &ConsulContext{}

	//TODO: deal addrs
	addrs := strings.Split(addrs_str, ",")
	for i, addr := range addrs {
		if !strings.Contains(addr, ":") {
			addrs[i] = addr + ":" + DEFAULT_CONSUL_API_PORT
		}
	}
	log.Debug("Consul context addrs is: ", addrs)
	cc.addrs     = addrs
	cc.index     = 0
	cc.addr_num  = len(addrs)
	cc.reg_loc   = reg_loc
	cc.dereg_loc = dereg_loc
	cc.read_loc  = read_loc
	cc.log       = log

	return cc, nil
}

func (cc *ConsulContext) getServer() string {
	//cc.index = (cc.index + 1) % cc.addr_num

	//return cc.addrs[cc.index]
	return cc.addrs[0]
}

func (cc *ConsulContext) OperateService(name, addr, id string, op int) (*ConsulResponse, error) {
	var resp *http.Response

	if op != READ && id == "" {
		cc.log.Error("Modify service without id")
		return nil, InternalServerError
	}

	r := &ConsulRequest{}

	r.Name = name
	r.Address = addr
	r.ID = id

	b, err := json.Marshal(r)
	if err != nil {
		cc.log.Error("Generate register json failed:", err)
		return nil, InternalServerError
	}
	cc.log.Debug(string(b))

	data := bytes.NewBuffer(b)

	switch op {
	case REGISTER:
		cc.log.Debug("http://" + cc.getServer() + cc.reg_loc, DEFAULT_CONTENT_HEADER, data)
		resp, err = http.Post("http://" + cc.getServer() + cc.reg_loc, DEFAULT_CONTENT_HEADER, data)
		break
	case DEREGISTER:
		cc.log.Debug("http://" + cc.getServer() + cc.dereg_loc + id)
		resp, err = http.Get("http://" + cc.getServer() + cc.dereg_loc + id)
		break
	case READ:
		resp, err = http.Get("http://" + cc.getServer() + cc.read_loc + name)
		break
	default: /* Should not reach here */
		cc.log.Error("Unknown operate code: ", op)
		return nil, InternalServerError
	}

	if err != nil {
		cc.log.Error("Opereate service to consul failed: ", err)
		return nil, BadGatewayError
	}

	defer resp.Body.Close()

	if resp.StatusCode != HTTP_OK {
		cc.log.Error("Opereate http status error: %d", resp.StatusCode)
		return nil, BadGatewayError
	}

	if op == READ {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			cc.log.Error("Read operate result body failed: ", err)
			return nil, InternalServerError
		}
		if string(body) == "[]" { //TODO: better way
			cc.log.Info("Service %s not found", name)
			return nil, NoContentError
		}
		cresp := &ConsulResponse{}
		err = json.Unmarshal(body, &cresp)
		if err != nil {
			cc.log.Error("Operate service decode json failed")
			return nil, InternalServerError
		}
		return cresp, nil
	}

	return nil, nil
}

func (cc *ConsulContext) RegisterService(name, addr string) error {
	id := name + "_" + fmt.Sprint(time.Now().Unix())
	_, err := cc.OperateService(name, addr, id, REGISTER)
	return err
}

func (cc *ConsulContext) DeRegisterService(name, addr string) error {
	found := true
	if addr != "" {
		found = false
	}

	resps, err := cc.OperateService(name, addr, "", READ)
	if err != nil {
		cc.log.Error("Deregister service failed")
		return err
	}
	if len(*resps) == 0 {
		cc.log.Info("Service %s not found", name)
		return NoContentError
	}
	for _, resp := range *resps {
		if addr != "" {
			if strings.EqualFold(resp.ServiceAddress, addr) {
				found = true
				_, err = cc.OperateService(name, addr, resp.ServiceID, DEREGISTER)
				if err != nil {
					return err
				}
			}
		} else {
			_, err = cc.OperateService(name, addr, resp.ServiceID, DEREGISTER)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		cc.log.Error("Not found service name: %s, addr: %s", name, addr)
		return NoContentError
	}
	return nil
}

func (cc *ConsulContext) ListService(name, addr string) (*ConsulResponse, error) {
	return cc.OperateService(name, addr, "", READ)
}

