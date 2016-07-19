package goblin

type GoblinRequest struct {
	Ip          string
	Uid         string
	Uuid        string
	Type        string
	Product     string
	Expire      int
	Action      string
	Ruleid      string
}

type GoblinReadResponse []GoblinRequest

type GoblinMessage struct {
	Startip     string
	Endip       string
	Uid         string
	Uuid        string
	Expire      string
	Punish      string
	Punish_args string
}

type ServerRequest struct {
	Addr    string
	Product string
}
type ServerResponse ServerRequest

