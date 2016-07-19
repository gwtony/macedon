package sample

import (
	//"io"
	//"strings"
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/test"
	//"net/http"
)

func TestInitHandler(t *testing.T) {
	log := test.TestInitlog()
	InitHandler("/", log)
	t.Log("init handler")
}

func TestRuleOperate(t *testing.T) {
	log := test.TestInitlog()
	h := InitHandler("/", log)
	h.RuleOperate("127.0.0.1:11111", "data", ADD_RULE)
	h.RuleOperate("127.0.0.1:11111", "data", DELETE_RULE)
	t.Log("rule operate done")
}

func TestRuleAdd(t *testing.T) {
	log := test.TestInitlog()
	h := InitHandler("/", log)
	h.RuleAdd("127.0.0.1", "nohost", 0, 1, "noheader")
	t.Log("rule add done")
}

func TestRuleDelete(t *testing.T) {
	log := test.TestInitlog()
	h := InitHandler("/", log)
	h.RuleDelete("127.0.0.1", "nohost", 0, 1, "noheader")
	t.Log("rule add done")
}

