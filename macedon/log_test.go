package macedon
import (
	"os"
	"fmt"
	"testing"
	"io/ioutil"
	"path/filepath"
)

var testTempDir string
func test_init_log() (* Log) {
    testTempDir, err := ioutil.TempDir("", "test_macedon_log")
    if err != nil {
        fmt.Errorf("tempDir: %v", err)
		return nil
    }

    path := filepath.Join(testTempDir, "test.log")
    log := GetLogger(path, "error")

	return log
}

func test_destroy_log() {
	os.RemoveAll(testTempDir)
}

func Test_Log(t *testing.T) {
	log := test_init_log()

	defer test_destroy_log()
	if log == nil {
		t.Fatal("init log failed")
	}
	t.Log("init log ok")
}
