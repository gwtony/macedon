package macedon

import (
	"time"
	"net/http"
	"io/ioutil"
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
	unreg_loc string
	read_loc  string

	log       *Log
}

func InitConsulContext(addrs, reg_loc, dereg_loc, read_loc string, log *Log) (*ConsulContext, error) {
	cc := &ConsulContext{}

	//TODO: deal addrs
	cc.reg_loc   = reg_loc
	cc.dereg_loc   = dereg_loc
	cc.read_loc = read_loc
	cc.log      = log

	return cc
}

func (cc *ConsulContext) getServer() (string, error) {
	cc.index = (cc.index + 1) % cc.addr_num

	return cc.addrs[cc.index], nil
}

func (cc *ConsulContext) OperateService(name, addr, id string, op int) (*ConsulResponse, error) {
	resp := nil

	if op == READ && id == "" {
		cc.log.Error("Modify service without id")
		return nil, InternalServerError
	}

	r := &ConsulRequest{}

	r.Name = name
	r.Address = addr

	b, err := json.Marshal(r)
	if err != nil {
		cc.log.Error("Generate register json failed:", err)
		return nil, InternalServerError
	}
	cc.log.Debug(string(b))

	data := bytes.NewBuffer([]byte(b))

	location := ""

	swithc op {
	case REG:
		location = cc.reg_loc
		break
	case UNREG:
		location = cc.dereg_loc
		break
	case READ:
		location = cc.read_loc
		break
	default: /* Should not reach here */
		cc.log.Error("Unknown operate code: ", op)
		return nil, InternalServerError
	}

	if op == UNREG {
		resp, err = http.Get(cc.getServer() + location + id)
	} else {
		resp, err = http.Post(cc.getServer() + location, DEFAULT_CONTENT_HEADER, data)
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

	resps, err := cc.OperateService(name, addr, READ)
	if err != nil {
		cc.log.Error("Deregister service failed")
		return err
	}
	if len(resps) == 0 {
		cc.log.Info("Service %s not found", name)
		return NoContentError
	}
	for _, resp := range resps {
		if addr != "" {
			if strings.EqualFold(resp.ServiceAddress, addr) {
				found = true
				_, err = cc.OperateService(name, addr, resp.ServiceID, UNREGISTER)
				if err != nil {
					return err
				}
				break
			}
		} else {
			_, err = cc.OperateService(name, addr, resp.ServiceID, UNREGISTER)
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

func (cc *ConsulContext) ListService(name string) (*ConsulResponse, error) {
	return cc.OperateService(name, addr, READ)
}

