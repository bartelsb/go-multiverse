package unixfs

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ipfs/go-merkledag/dagutils"
	ufs "github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"
)

func TestAddFile(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	original, err := ioutil.ReadFile("testdata/a.txt")
	if err != nil {
		t.Fatal("failed to read file")
	}

	node, err := Add(ctx, dag, "testdata/a.txt", nil)
	if err != nil {
		t.Fatal("failed to add file")
	}

	r, err := ufsio.NewDagReader(ctx, node, dag)
	if err != nil {
		t.Fatal("failed to read node")
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal("failed to read node")
	}

	if !bytes.Equal(original, data) {
		t.Error("file data does not match")
	}
}

func TestAddDir(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	node, err := Add(ctx, dag, "testdata", nil)
	if err != nil {
		t.Fatal("failed to add")
	}

	ufsdir, err := ufsio.NewDirectoryFromNode(dag, node)
	if err != nil {
		t.Fatal("failed to read node")
	}

	if _, err := ufsdir.Find(ctx, "b"); err != nil {
		t.Error("failed to find file")
	}

	if _, err := ufsdir.Find(ctx, "l"); err != nil {
		t.Error("failed to find file")
	}

	if _, err := ufsdir.Find(ctx, "a.txt"); err != nil {
		t.Error("failed to find file")
	}

	if _, err := ufsdir.Find(ctx, "b.txt"); err != nil {
		t.Error("failed to find file")
	}

	if _, err := ufsdir.Find(ctx, "o.txt"); err != nil {
		t.Error("failed to find file")
	}

	if _, err := ufsdir.Find(ctx, "r.txt"); err != nil {
		t.Error("failed to find file")
	}
}

func TestAddSymlink(t *testing.T) {
	ctx := context.Background()
	dag := dagutils.NewMemoryDagService()

	node, err := Add(ctx, dag, "testdata/l", nil)
	if err != nil {
		t.Fatal("failed to add")
	}

	fsnode, err := ufs.ExtractFSNode(node)
	if err != nil {
		t.Fatal("failed to extract fsnode")
	}

	target, err := os.Readlink("testdata/l")
	if err != nil {
		t.Fatal("failed to read link")
	}

	if target != string(fsnode.Data()) {
		t.Error("unexpected symlink data")
	}
}

func TestIgnore(t *testing.T) {
	ignore := Ignore{
		"*.exe",
		"baz/*",
	}

	if !ignore.Match("foo.exe") {
		t.Error("expected ignore to match")
	}

	if !ignore.Match("foo/bar.exe") {
		t.Error("expected ignore to match")
	}

	if !ignore.Match("baz/bar") {
		t.Error("expected ignore to match")
	}
}
