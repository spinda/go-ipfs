//go:generate go-bindata -pkg=assets init-doc ../vendor/dir-index-html-v1.0.0
//go:generate gofmt -w bindata.go

package assets

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ipfs/go-ipfs/blocks/key"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	uio "github.com/ipfs/go-ipfs/unixfs/io"
)

// initDocPaths lists the paths for the docs we want to seed during --init
var initDocPaths = []string{
	filepath.Join("init-doc", "about"),
	filepath.Join("init-doc", "readme"),
	filepath.Join("init-doc", "help"),
	filepath.Join("init-doc", "contact"),
	filepath.Join("init-doc", "security-notes"),
	filepath.Join("init-doc", "quick-start"),
}

// SeedInitDocs adds the list of embedded init documentation to the passed node, pins it and returns the root key
func SeedInitDocs(nd *core.IpfsNode) (*key.Key, error) {
	return addAssetList(nd, initDocPaths)
}

var initDirIndex = []string{
	filepath.Join("..", "vendor", "dir-index-html-v1.0.0", "knownIcons.txt"),
	filepath.Join("..", "vendor", "dir-index-html-v1.0.0", "dir-index.html"),
}

func SeedInitDirIndex(nd *core.IpfsNode) (*key.Key, error) {
	return addAssetList(nd, initDirIndex)
}

func addAssetList(nd *core.IpfsNode, l []string) (*key.Key, error) {
	dirb := uio.NewDirectory(nd.DAG)

	for _, p := range l {
		d, err := Asset(p)
		if err != nil {
			return nil, fmt.Errorf("assets: could load Asset '%s': %s", p, err)
		}

		s, err := coreunix.Add(nd, bytes.NewBuffer(d))
		if err != nil {
			return nil, fmt.Errorf("assets: could not Add '%s': %s", p, err)
		}

		fname := filepath.Base(p)
		k := key.B58KeyDecode(s)
		if err := dirb.AddChild(nd.Context(), fname, k); err != nil {
			return nil, fmt.Errorf("assets: could not add '%s' as a child: %s", fname, err)
		}
	}

	dir := dirb.GetNode()
	dkey, err := nd.DAG.Add(dir)
	if err != nil {
		return nil, fmt.Errorf("assets: DAG.Add(dir) failed: %s", err)
	}

	if err := nd.Pinning.Pin(nd.Context(), dir, true); err != nil {
		return nil, fmt.Errorf("assets: Pinning on init-docu failed: %s", err)
	}

	if err := nd.Pinning.Flush(); err != nil {
		return nil, fmt.Errorf("assets: Pinning flush failed: %s", err)
	}

	return &dkey, nil
}
