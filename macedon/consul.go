package macedon

import (
	"net/http"
	"io/ioutil"
)

var {
	DEFAULT_HEADER = "application/json;charset=utf-8"
	REGISTER_LOC   = "/v1/agent/register"
	UNREGISTER_LOC = "/v1/agent/unregister"
	READ_LOC       = "/v1/agent/" //TODO:
}

const (
	REG = iota
	UNREG
	READ
)

type ConsulContext struct {
	addrs     []string
	index     int
	addr_num  int

	domain    string
	reg_loc   string
	unreg_loc string
	read_loc  string

	log       *Log
}

func InitConsulContext(addrs, port, domain, op_loc, read_loc string, log *Log) (*ConsulContext, error) {
	cc := &ConsulContext{}

	//TODO: deal addrs
	cc.domain   = domain
	cc.op_loc   = op_loc
	cc.read_loc = read_loc
	cc.log      = log

	return cc
}

func (cc *ConsulContext) getServer() (string, error) {
	cc.index = (cc.index + 1) % cc.addr_num

	return cc.addrs[cc.index], nil
}

func (cc *ConsulContext) OperateService(name, addr string, op int) error {
	r := &ConsulRequest{}

	r.Name = name
	r.Address = addr

	b, err := json.Marshal(r)
	if err != nil {
		cc.log.Error("Generate register json failed:", err)
		return InternalServerError
	}
	cc.log.Debug(string(b))

	data := bytes.NewBuffer([]byte(b))

	location := ""
	swithc op {
	case REG:
		location = REGISTER_LOC
		break
	case UNREG:
		location = UNREGISTER_LOC
		break
	case READ:
		location = READ_LOC
		break
	default: /* Should not reach here */
		cc.log.Error("Unknown operate code: ", op)
		return nil
	}

	resp, err := http.Post(cc.getServer() + location, DEFAULT_HEADER, data)
	if err != nil {
		cc.log.Error("Opereate service to consul failed: ", err)
		return BadGatewayError
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cc.log.Error("Read operate result body failed: ", err)
		return InternalServerError
	}

	if body.StatusCode != HTTP_OK {
		cc.log.Error("Opereate http status error: %d", body.StatusCode)
		return BadGatewayError
	}

	return nil
}

func (cc *ConsulContext) RegisterService(name, addr string) error {
	return cc.OperateService(name, addr, REG)
}
func (cc *ConsulContext) UnRegisterService(name string) error {
	return cc.OperateService(name, addr, UNREG)
}

func (cc *ConsulContext) ListService(name string) error {
	return cc.OperateService(name, addr, READ)
}

