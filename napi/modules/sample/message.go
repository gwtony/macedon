package sample

type MatchRule struct {
	Band   []int
	Host   []string
	Expire int
}

type ActionRule struct {
	Type  string
	Value string
}

type SampleRequest struct {
	Type    string
	Match   MatchRule
	Action  ActionRule
	Ruleid  string
	Host    string
	Product string
}

type SampleRuleResponse []SampleRule
type SampleRule SampleRequest

type SampleAddResponse struct {
	Ruleid string
}

type ServerRequest struct {
	Addr    string
	Product string
}
type ServerResponse ServerRequest

type ClearRequest struct {
	Expire int
	Token  string
}
