package macedon

import (
	"os"
	//"fmt"
	"testing"
	"io/ioutil"
	"path/filepath"
)
var testKeyPath, testKeyTempDir string

func testGenerateKey() error {
	testKeyTempDir, err := ioutil.TempDir("", "test_server")
	if err != nil {
		return err
	}

	testKeyPath = filepath.Join(testKeyTempDir, "test_id_rsa")
	err = ioutil.WriteFile(testKeyPath, []byte(test_key), 0644)
	if err != nil {
		return err
	}

	return nil
}

func testRemoveKey() {
	if testKeyTempDir != "" {
		os.RemoveAll(testKeyTempDir)
	//} else {
	//	fmt.Println("key is null")
	}
}

func TestServerNotPurgeOk(t *testing.T) {
	conf := &Config{
				"192.168.0.1",
				"/test",
				"22", "root", "./id_isa.pub", 10,
				"127.0.0.1", "echo", 0,
				"192.168.0.1", "", "", "", "test.com",
				"",
				"",
	}

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	_, err := InitServer(conf, log)
	if err != nil {
		t.Fatal("Test server without purge failed!")
	}
	t.Log("Test server without purge ok.")
}

func TestServerPurgeOk(t *testing.T) {
	err := testGenerateKey()
	if err != nil {
		t.Fatal("generate key failed")
	}

	conf := &Config{
				"172.30.19.33",
				"/test2",
				"22", "root", testKeyPath, 10,
				"127.0.0.1", "echo", 1,
				"172.30.19.33", "", "", "", "test.com",
				"",
				"",
	}

	log := testInitlog()
	if log == nil {
		t.Fatal("init log failed")
	}
	defer testDestroylog()

	_, err = InitServer(conf, log)
	if err != nil {
		t.Fatal("Test server with purge failed!")
	}
	testRemoveKey()
	t.Log("Test server with purge ok.")
}
