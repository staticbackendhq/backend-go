package backend_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/staticbackendhq/backend-go"
)

func TestUploadFile(t *testing.T) {
	f, err := os.Open("storage_test.go")
	if err != nil {
		t.Fatal(err)
	}

	res, err := backend.StoreFile(token, "unittest.go", f)
	if err != nil {
		t.Error(err)
	} else if !strings.HasPrefix(res.URL, "http://localhost:8099/tmp") {
		t.Errorf("expected URL to have http as prefix got %s", res.URL)
	}

	orig, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := backend.DownloadFile(token, res.URL)
	if err != nil {
		t.Error(err)
	} else if len(orig) != len(buf) {
		t.Errorf("uploaded buffer and disk file differ, expected length %d got %d", len(orig), len(buf))
	}
}
