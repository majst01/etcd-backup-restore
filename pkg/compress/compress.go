package compress

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/gardener/etcd-backup-restore/pkg/snapstore"
	"github.com/mholt/archiver/v3"
)

type (
	// Compressor compress/decompress backup data before/after sending/receiving from storage
	Compressor struct {
		extension string
		enabled   bool
	}
)

// New Returns a new Compressor
func New(method string) (*Compressor, error) {
	c := &Compressor{
		enabled: true,
	}
	switch method {
	case "", "none":
		c.extension = ""
		c.enabled = false
	case "tar":
		c.extension = ".tar"
	case "auto", "targz":
		c.extension = ".tar.gz"
	case "tarlz4":
		c.extension = ".tar.lz4"
	default:
		return nil, fmt.Errorf("unsupported compression method: %s", method)
	}
	return c, nil
}

// Compress the given snapshot, SnapName is adjusted matching compression algorithm
func (c *Compressor) Compress(snap *snapstore.Snapshot) error {
	if !c.enabled {
		return nil
	}
	err := archiver.Archive([]string{snap.SnapDir}, path.Join(snap.SnapDir, snap.SnapName))
	if err != nil {
		return err
	}
	snap.SnapName = snap.SnapName + c.extension
	return nil
}

// Decompress the given snapshot, if snapshot is not compressed, nothing happens
func (c *Compressor) Decompress(snap *snapstore.Snapshot) error {
	if !c.enabled {
		return nil
	}
	// skip decompression for snaps stored without compression
	if !isCompressed(snap) {
		return nil
	}
	err := archiver.Unarchive(path.Join(snap.SnapDir, snap.SnapName), filepath.Dir(snap.SnapDir))
	if err != nil {
		return err
	}
	snap.SnapName = strings.TrimSuffix(snap.SnapName, c.extension)
	return nil
}

func isCompressed(snap *snapstore.Snapshot) bool {
	compressed := false
	for _, ext := range []string{".tar", ".tar.gz", ".tar.lz4"} {
		if strings.HasSuffix(snap.SnapName, ext) {
			compressed = true
		}
	}
	return compressed
}
