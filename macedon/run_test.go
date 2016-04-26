package macedon

import (
	"os"
	//"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const test_macedon_run_conf string = `
[default]
addr: a.b.c.d

log: macedon.log
level: debug

location: /dns
purge_ips: 127.0.0.1
purge_cmd: touch /tmp/ggg01
ssh_key: /root/.ssh/id_rsa
ssh_port: 22
ssh_user: root
ssh_timeout: 20


consul_addrs: 127.0.0.1:8500
register_location: /v1/agent/service/register
deregister_location: /v1/agent/service/deregister/
read_location: /v1/catalog/service/
domain: test.com
`

func TestParseOptionVersion(t *testing.T) {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-v")
	} else {
		os.Args[1] = "-v"
	}
	ParseOption()
	t.Log("parse version done")
}
func TestParseOptionHelp(t *testing.T) {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-h")
	} else {
		os.Args[1] = "-h"
	}
	ParseOption()
	t.Log("parse help done")
}
func TestParseOptionConfig(t *testing.T) {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-f")
		os.Args = append(os.Args, "a")
	} else if len(os.Args) == 2{
		os.Args[1] = "-f"
		os.Args = append(os.Args, "a")
	} else {
		os.Args[1] = "-f"
		os.Args[2] = "a"
	}
	ParseOption()
	t.Log("parse config done")
}
func TestParseOptionInvalid(t *testing.T) {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-a")
	} else {
		os.Args[1] = "-a"
	}
	ParseOption()
	t.Log("parse invalid done")
}


func TestRun(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_log")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	path := filepath.Join(tempDir, "test_conf")
	err = ioutil.WriteFile(path, []byte(test_macedon_run_conf), 0644)
	if err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	switch len(os.Args) {
	case 1:
		os.Args = append(os.Args, "-f")
		os.Args = append(os.Args, path)
		break
	case 2:
		os.Args[1] = "-f"
		os.Args = append(os.Args, path)
		break
	default:
		os.Args[1] = "-f"
		os.Args[2] = path
		break
	}
	Run()
}
