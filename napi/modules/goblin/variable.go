package goblin

const (
	ADD_RULE                = iota
	DELETE_RULE
	READ_RULE

	GOBLIN_RULE_ADD_LOCATION       = "/rule/add"
	GOBLIN_RULE_DELETE_LOCATION    = "/rule/delete"
	GOBLIN_RULE_READ_LOCATION      = "/rule/read"

	GOBLIN_SERVER_ADD_LOCATION     = "/server/add"
	GOBLIN_SERVER_DELETE_LOCATION  = "/server/delete"
	GOBLIN_SERVER_READ_LOCATION    = "/server/read"

	GOBLIN_ADD_LOCATION            = "/add"
	GOBLIN_DELETE_LOCATION         = "/del"

	GOBLIN_LOCATION                = "/goblin"
	GOBLIN_API_LOCATION            = "/goblin"

	EMPTY_IP                       = "0.0.0.0"
	EMPTY_UUID                     = "0"
	EMPTY_UID                      = "0"
)
