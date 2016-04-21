package macedon

import (
	"time"
	"testing"
)

const test_key string = `
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDW6jay6kk69YiCnAiKzVxwD+kgelOMEgAa2UNLuEwEMMN5WXYd1uZKQ2Wzju9FeIWeKXwITWk+K/y63LeUAJG1E0IBb+/OsdR/oh9VcYZEvGrGgPyuIHYFpF+2mm9vJby91xPT4z72Af/6BCQ0XL1OFQ4PV2EMCuNjE6IenGk1wuq2P4xaD3rWghM4m4w7TN5SJMuwtGXchYwZv3wKI70VlrEc/4FO1LqxMMJ8v350UNH3MWW5Nsq0N1kGujFjPnoC1LC97lMRtAWeiLAVyRPRlMHk7bSlXnuG7GnqkDMqzL6eARZvqWYEK7Ulz1ZMWHgjJzuwQj1FMouQcQDVWeRj root@cent7-dev
`

func Test_InitSshContext_OK(t *testing.T) {
	log := test_init_log()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer test_destroy_log()

	InitSshContext(test_key, "root", time.Second * 5, log)
	//sc, err := InitSshContext(test_key, "root", time.Second * 5, log)
	//if err != nil {
	//	//t.Fatal("init ssh context failed")
	//	t.Log("init ssh context ok")
	//}
	//t.Log("init ssh context ok")
}
