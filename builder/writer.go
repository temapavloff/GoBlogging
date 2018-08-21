package builder

import (
	"GoBlogging/config"
	"GoBlogging/layout"
	"GoBlogging/pages"
	"os"
	"path"
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

// Prepare - cleans the output directory
func (w *Writer) Prepare() error {
	outDir := w.config.GetAbsPath(w.config.Output)

	if err := os.RemoveAll(outDir); err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	assetsPath := w.layout.GetAssetsPath()
	if assetsPath != "" {
		copyAll(assetsPath, outDir)
	}

	return nil
}

// Write - writes page to disk
func (w *Writer) Write(nodeCh <-chan pages.Node,
	errCh chan<- error, wg *sync.WaitGroup) {
	var err error
	for n := range nodeCh {
		switch n.Type {
		case pages.IndexType:
			err = w.writeIndex(n.Index)
		case pages.TagType:
			err = w.writeTag(n.Tag)
		case pages.PostType:
			err = w.writePost(n.Post)
		}
		if err != nil {
			errCh <- err
		}
		wg.Done()
	}
}

func (w *Writer) writeIndex(index *pages.Index) error {
	f, err := os.Create(path.Join(index.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	return w.layout.RenderIndex(f, index)
}

func (w *Writer) writeTag(tag *pages.Tag) error {
	if err := os.MkdirAll(tag.Output, 0755); err != nil {
		return err
	}

	f, err := os.Create(path.Join(tag.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	return w.layout.RenderTag(f, tag)
}

func (w *Writer) writePost(post *pages.Post) error {
	if err := os.MkdirAll(post.OutputPath, 0755); err != nil {
		return err
	}

	if err := copyExclude(post.InputPath, post.OutputPath, ".md"); err != nil {
		return err
	}

	f, err := os.Create(path.Join(post.OutputPath, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	return w.layout.RenderPost(f, post)
}
