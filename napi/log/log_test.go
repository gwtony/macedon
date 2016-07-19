package log
import (
	"os"
	"testing"
	"io/ioutil"
	"path/filepath"
)

func TestLog(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_napi_log")
	if err != nil {
		t.Fatal("create temp dir failed")
	}

	path := filepath.Join(tempDir, "test.log")
	log := GetLogger(path, "debug")

	if log == nil {
		t.Fatal("init log failed")
	}
	os.RemoveAll(tempDir)
	t.Log("init log ok")
}
