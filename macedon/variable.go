package macedon

const (
	REGISTER = iota
	DEREGISTER
	READ

	DEFAULT_SUB_ZONE = ".service."
	DEFAULT_LOG_PATH = "../log/macedon.log"

	DEFAULT_CONTENT_HEADER = "application/json;charset=utf-8"

	DEFAULT_REGISTER_LOC   = "/v1/agent/service/register"
	DEFAULT_DEREGISTER_LOC = "/v1/agent/service/deregister/"
	DEFAULT_READ_LOC       = "/v1/catalog/service/"

	DefaultScpCmd = "/usr/bin/scp -qrt "

	DEFAULT_CONF_PATH         = "../conf"
	DEFAULT_CONF              = "macedon.conf"

	DEFAULT_CREATE_LOCATION   = "/create"
	DEFAULT_DELETE_LOCATION   = "/delete"
	DEFAULT_READ_LOCATION     = "/read"

	DEFAULT_CONSUL_API_PORT   = "8600"

	DEFAULT_SSH_TIMEOUT       = 5

	HTTP_OK   = 200

)
