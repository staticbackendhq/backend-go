package backend_test

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/staticbackendhq/backend-go"
)

func TestUploadFile(t *testing.T) {
	f, err := os.Open("storage_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res, err := backend.StoreFile(token, "unittest.go", f)
	if err != nil {
		t.Error(err)
	} else if !strings.HasPrefix(res.URL, "http://localhost:8099/localfs") {
		t.Errorf("expected URL to have http as prefix got %s", res.URL)
	}

	orig, err := os.ReadFile("storage_test.go")
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

func TestStorageUsage(t *testing.T) {
	testToken := newStorageTestAccount(t)

	before, err := backend.StorageUsage(testToken)
	if err != nil {
		t.Fatal(err)
	}

	size1 := int64(1024)
	size2 := int64(2048)
	uploadStorageFixture(t, testToken, "usage-1.txt", size1)
	uploadStorageFixture(t, testToken, "usage-2.txt", size2)

	after, err := backend.StorageUsage(testToken)
	if err != nil {
		t.Fatal(err)
	}

	if after.Bytes-before.Bytes != size1+size2 {
		t.Fatalf("expected usage delta %d bytes got %d", size1+size2, after.Bytes-before.Bytes)
	}

	wantGB := math.Round((float64(after.Bytes)/1_000_000_000)*100) / 100
	if after.GB != wantGB {
		t.Fatalf("expected usage %v GB got %v", wantGB, after.GB)
	}
}

func TestListFiles(t *testing.T) {
	testToken := newStorageTestAccount(t)

	for i := 0; i < 30; i++ {
		uploadStorageFixture(t, testToken, fmt.Sprintf("list-%02d.txt", i), int64(i+1))
	}

	result, err := backend.ListFiles(testToken, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result.Page != 1 || result.Size != 25 {
		t.Fatalf("expected page 1 size 25 got page=%d size=%d", result.Page, result.Size)
	}
	if result.Total != 30 {
		t.Fatalf("expected 30 files got %d", result.Total)
	}
	if len(result.Results) != 25 {
		t.Fatalf("expected 25 results got %d", len(result.Results))
	}
	if result.Results[0].Size != 30 {
		t.Fatalf("expected latest upload first with size 30 got %d", result.Results[0].Size)
	}

	result, err = backend.ListFiles(testToken, &backend.ListFilesParams{Page: 2})
	if err != nil {
		t.Fatal(err)
	}

	if result.Page != 2 {
		t.Fatalf("expected page 2 got %d", result.Page)
	}
	if len(result.Results) != 5 {
		t.Fatalf("expected 5 results on page 2 got %d", len(result.Results))
	}
	if result.Results[0].Size != 5 {
		t.Fatalf("expected page 2 to start at size 5 got %d", result.Results[0].Size)
	}

	result, err = backend.ListFiles(testToken, &backend.ListFilesParams{SortBy: "size"})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Results) == 0 {
		t.Fatal("expected non-empty results")
	}
	if result.Results[0].Size != 30 {
		t.Fatalf("expected largest file first got %d", result.Results[0].Size)
	}
}

func newStorageTestAccount(t *testing.T) string {
	t.Helper()

	email := fmt.Sprintf("storage-%d@test.com", time.Now().UnixNano())
	testToken, err := backend.Register(email, "devpw1234")
	if err != nil {
		t.Fatal(err)
	}

	return testToken
}

func uploadStorageFixture(t *testing.T, testToken, filename string, size int64) backend.StoreFileResult {
	t.Helper()

	body := bytes.NewReader(bytes.Repeat([]byte("a"), int(size)))
	res, err := backend.StoreFile(testToken, filename, body)
	if err != nil {
		t.Fatal(err)
	}

	return res
}
