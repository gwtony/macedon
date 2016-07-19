package config

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/test"
)


func testReadConf(t *testing.T, data string) *Config {
	conf := &Config{}
	tempDir, err := ioutil.TempDir("", "test_log")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	path := filepath.Join(tempDir, "test_conf")
	err = ioutil.WriteFile(path, []byte(test.Test_napi_conf), 0644)
	if err != nil {
		t.Fatalf("writeFile: %v", err)
	}

	err = conf.ReadConf(path)
	if err != nil {
		t.Fatal("Test read conf failed")
	}
	t.Log("Test read conf ok")
	return conf
}

func TestReadConfOk(t *testing.T) {
	testReadConf(t, test.Test_napi_conf)
}

func TestParseConfOk(t *testing.T) {
	c := testReadConf(t, test.Test_napi_conf)
	c.ParseConf()
}
