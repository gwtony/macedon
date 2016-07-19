package sample

const (
	ADD_RULE                = iota
	DELETE_RULE
	READ_RULE


	EQUAL_HOST                    = "equal"
	LIKE_HOST                     = "like"
	SAMPLE_TYPE                   = "sample"
	SAMPLE_ACTION_TYPE            = "insert_header"

	SAMPLE_RULE_ADD_LOCATION      = "/rule/add"
	SAMPLE_RULE_DELETE_LOCATION   = "/rule/delete"
	SAMPLE_RULE_READ_LOCATION     = "/rule/read"
	SAMPLE_RULE_CLEAR_LOCATION    = "/rule/clear"

	SAMPLE_SERVER_ADD_LOCATION    = "/server/add"
	SAMPLE_SERVER_DELETE_LOCATION = "/server/delete"
	SAMPLE_SERVER_READ_LOCATION   = "/server/read"

	SAMPLE_ADD_LOCATION           = "/add"
	SAMPLE_DELETE_LOCATION        = "/del"

	SAMPLE_LOCATION               = "/sample"
	SAMPLE_API_LOCATION           = "/sample"

	EMPTY_IP                      = "0.0.0.0"
	EMPTY_UUID                    = "0"
	EMPTY_UID                     = "0"

	SAMPLE_TOKEN                  = "secret"
)
