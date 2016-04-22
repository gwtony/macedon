package macedon
import (
	"os"
	"fmt"
	"testing"
	"io/ioutil"
	"path/filepath"
)

var testTempDir string
func testInitlog() (* Log) {
    testTempDir, err := ioutil.TempDir("", "test_macedon_log")
    if err != nil {
        fmt.Errorf("tempDir: %v", err)
		return nil
    }

    path := filepath.Join(testTempDir, "test.log")
    log := GetLogger(path, "debug")

	return log
}

func testDestroylog() {
	if testTempDir == "" {
		os.RemoveAll(testTempDir)
	} else {
		fmt.Println("temp dir is null")
	}
}

func TestLog(t *testing.T) {
	log := testInitlog()

	defer testDestroylog()
	if log == nil {
		t.Fatal("init log failed")
	}
	t.Log("init log ok")
}
