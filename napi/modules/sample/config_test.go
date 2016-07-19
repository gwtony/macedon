package sample

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"testing"
	"git.lianjia.com/lianjia-sysop/napi/config"
)


const Test_sample_conf string = `
[sample]
addr: 172.30.23.39:12345

log: napi.log
level: debug

location: /sample
api_location: /sample_api
mysql_addr: 127.0.0.1:3306
mysql_dbname: test
mysql_dbuser: test
mysql_dbpwd: test
`
var Test_sample_conf_file = "sample_conf_file"

func testSampleReadConf(t *testing.T) {
	conf := &config.Config{}
	tempDir, err := ioutil.TempDir("", "test_log")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	path := filepath.Join(tempDir, "test_conf")
	err = ioutil.WriteFile(path, []byte(Test_sample_conf_file), 0644)
	if err != nil {
		t.Fatalf("writeFile: %v", err)
	}

	err = conf.ReadConf(path)
	if err != nil {
		t.Fatal("Test read conf failed")
	}
	sconf := SampleConfig{}
	err = sconf.ParseConfig(conf)
	if err != nil {
		t.Fatal("Parse sample config failed")
	}
	t.Log("Test read conf done")
}
