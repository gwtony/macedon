package macedon

import (
	"fmt"
	"time"
	"bytes"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

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
	var addrs []string

	taddrs := strings.Split(addrs_str, ",")
	for _, addr := range taddrs {
		if addr == "" {
			continue
		}
		if !strings.Contains(addr, ":") {
			addrs = append(addrs, addr + ":" + DEFAULT_CONSUL_API_PORT)
		} else {
			addrs = append(addrs, addr)
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
	cc.index = (cc.index + 1) % cc.addr_num

	return cc.addrs[cc.index]
	//return cc.addrs[0]
}

func (cc *ConsulContext) OperateService(name, addr, id, tags string, op int) (*ConsulResponse, error) {
	var resp *http.Response

	if op != READ && id == "" {
		cc.log.Error("Modify service without id")
		return nil, InternalServerError
	}

	r := &ConsulRequest{}

	r.Name = name
	r.Address = addr
	r.ID = id
	if tags != "" {
		r.Tags = append(r.Tags, tags)
	}

	b, err := json.Marshal(r)
	if err != nil {
		cc.log.Error("Generate register json failed:", err)
		return nil, InternalServerError
	}
	cc.log.Debug(string(b))

	data := bytes.NewBuffer(b)

	switch op {
	case REGISTER:
		resp, err = http.Post("http://" + cc.getServer() + cc.reg_loc, DEFAULT_CONTENT_HEADER, data)
		break
	case DEREGISTER:
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

		if len(body) == 0 ||
			strings.EqualFold(string(body), "[]") ||
			strings.EqualFold(string(body), "[]\n") { //TODO: better way
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

func (cc *ConsulContext) RegisterService(name, addr, tags string) error {
	id := name + "_" + fmt.Sprint(time.Now().Unix()) + "_" + strconv.Itoa(rand.Intn(10000))
	_, err := cc.OperateService(name, addr, id, tags, REGISTER)
	return err
}

func (cc *ConsulContext) DeRegisterService(name, addr, tags string) error {
	found := true
	if addr != "" {
		found = false
	}

	resps, err := cc.OperateService(name, addr, "", tags, READ)
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
				_, err = cc.OperateService(name, addr, resp.ServiceID, tags, DEREGISTER)
				if err != nil {
					return err
				}
			}
		} else {
			_, err = cc.OperateService(name, addr, resp.ServiceID, tags, DEREGISTER)
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

func (cc *ConsulContext) ListService(name, addr, tags string) (*ConsulResponse, error) {
	found := false

	cc.log.Debug("In list service")
	resps, err := cc.OperateService(name, addr, "", tags, READ)
	if err != nil {
		cc.log.Error("list service failed")
		return nil, err
	}

	if len(*resps) == 0 {
		cc.log.Info("Service %s not found", name)
		return nil, NoContentError
	}

	cc.log.Debug("Deal operate service result")

	if addr == "" && tags == "" {
		return resps, err
	}
	ret := &ConsulResponse{}
	for _, resp := range *resps {
		amatch := false
		tmatch := false
		if addr != "" {
			if strings.EqualFold(resp.ServiceAddress, addr) {
				amatch = true
			}
		} else {
			amatch = true
		}
		if tags != "" {
			for _, stag := range resp.ServiceTags {
				if strings.EqualFold(stag, tags) {
					tmatch = true
				}
			}
		} else {
			tmatch = true
		}
		if amatch && tmatch {
			*ret = append(*ret, resp)
			found = true
		}
	}

	if !found {
		cc.log.Error("Not found service name: %s, addr: %s", name, addr)
		return nil, NoContentError
	}

	return ret, nil
}
