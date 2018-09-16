package builder

import (
	"GoBlogging/config"
	"GoBlogging/helpers"
	"GoBlogging/layout"
	"GoBlogging/pages"
	"os"
	"path/filepath"
	"sync"
)

// Writer - writes down pages to disk
type Writer struct {
	config *config.Config
	layout layout.Layout
}

// NewWriter - creates new writer instance
func NewWriter(c *config.Config, l layout.Layout) *Writer {
	return &Writer{
		config: c,
		layout: l,
	}
}

func cleanup(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		// Skip dotnames to keep .git directory
		// Also keep CHANE file to support github pages
		if name[0] == '.' || name == "CNAME" {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Prepare - cleans the output directory
func (w *Writer) Prepare() error {
	outDir := w.config.GetAbsPath(w.config.Output)

	if err := cleanup(outDir); err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	assetsPath := w.layout.GetAssetsPath()

	// Its OK if template doesn't have assets
	if _, err := os.Stat(assetsPath); err == nil && assetsPath != "" {
		return helpers.CopyAll(assetsPath, outDir+"/assets")
	}

	return nil
}

// Write - writes page to disk
func (w *Writer) Write(pageCh <-chan pages.Page,
	errCh chan<- error, wg *sync.WaitGroup) {
	for n := range pageCh {
		if err := n.Write(w.layout); err != nil {
			errCh <- err
		}
		wg.Done()
	}
}
