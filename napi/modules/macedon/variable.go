package macedon

const (
	VERSION                 = "0.1 alpha"

	ADD                     = iota
	DELETE
	READ

	API_CONTENT_HEADER      = "application/json;charset=utf-8"
	ETCD_CONTENT_HEADER     = "application/x-www-form-urlencoded"

	ADD_METHOD              = "PUT"
	DELETE_METHOD           = "DELETE"
	CONTENT_HEADER          = "Content-Type"

	MACEDON_CREATE_LOC      = "/add"
	MACEDON_DELETE_LOC      = "/delete"
	MACEDON_READ_LOC        = "/read"
	MACEDON_ADD_SERVER_LOC  = "/server/add"
	MACEDON_DEL_SERVER_LOC  = "/server/delete"
	MACEDON_READ_SERVER_LOC = "/server/read"

	DEFAULT_ETCD_PORT       = "2379"
	DEFAULT_SKYDNS_LOC      = "/v2/keys/skydns/"
	DEFAULT_TTL             = 60

	DEFAULT_TIMEOUT     = 5

)
