package handler

const (
	// VERSION version
	VERSION                   = "0.1 alpha"

	//ADD operation in handler
	ADD                       = iota
	DELETE
	UPDATE
	READ

	API_CONTENT_HEADER        = "application/json;charset=utf-8"
	ETCD_CONTENT_HEADER       = "application/x-www-form-urlencoded"

	ADD_METHOD                = "PUT"
	DELETE_METHOD             = "DELETE"
	CONTENT_HEADER            = "Content-Type"

	MACEDON_ADD_LOC           = "/add"
	MACEDON_DELETE_LOC        = "/delete"
	MACEDON_UPDATE_LOC        = "/update"
	MACEDON_READ_LOC          = "/read"
	MACEDON_SCAN_LOC          = "/scan"
	MACEDON_ADD_SERVER_LOC    = "/server/add"
	MACEDON_DELETE_SERVER_LOC = "/server/delete"
	MACEDON_READ_SERVER_LOC   = "/server/read"

	MACEDON_LOC               = "/macedon"
	DEFAULT_ETCD_MACEDON_LOC  = "/v2/keys/macedon/node/purge/"
	DEFAULT_ETCD_PORT         = "2379"
	DEFAULT_ETCD_TIMEOUT      = 5
	DEFAULT_SKYDNS_LOC        = "/v2/keys/skydns/"
	DEFAULT_ARPA_LOC          = "/v2/keys/skydns/arpa/in-addr/"
	DEFAULT_SCAN_ARGS         = "/?recursive=true"
	DEFAULT_ETCD_CAS          = "?prevValue="
	DEFAULT_TTL               = 60

	DEFAULT_PURGE_CMD         = "REQ: PurgeCacheEntry"
	DEFAULT_PURGE_PORT        = "9180"
	DEFAULT_PURGE_TIMEOUT     = 5
	DEFAULT_PURGE_SERVER_LOC  = "/v2/keys/macedon/node/purge/"

	DEFAULT_TRIM_KEY          = "/skydns/"
	DEFAULT_TRIM_ARPA_KEY     = DEFAULT_TRIM_KEY + "arpa/in-addr/"
	DEFAULT_TRIM_SERVER_KEY   = "/macedon/node/purge/"

	DEFAULT_TIMEOUT           = 5

	MACEDON_TOKEN             = "macedon_token"
)
