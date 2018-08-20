package builder

import (
	"GoBlogging/config"
	"GoBlogging/layout"
	"GoBlogging/pages"
	"os"
	"path"
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

// GetWriterFn - returns Writer function
func (w *Writer) GetWriterFn(doneCh chan<- bool, errCh chan<- error) pages.PageWalker {
	return func(n pages.Node) error {
		switch n.Type {
		case pages.IndexType:
			go w.writeIndex(n.Index, doneCh, errCh)
		case pages.TagType:
			go w.writeTag(n.Tag, doneCh, errCh)
		case pages.PostType:
			go w.writePost(n.Post, doneCh, errCh)
		}
		return nil
	}
}

func (w *Writer) writeIndex(index *pages.Index, doneCh chan<- bool, errCh chan<- error) {
	f, err := os.Create(path.Join(index.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		errCh <- err
		return
	}

	if err := f.Chmod(0644); err != nil {
		errCh <- err
		return
	}

	if err := w.layout.RenderIndex(f, index); err != nil {
		errCh <- err
		return
	}

	doneCh <- true
}

func (w *Writer) writeTag(tag *pages.Tag, doneCh chan<- bool, errCh chan<- error) {
	if err := os.MkdirAll(tag.Output, 0755); err != nil {
		errCh <- err
		return
	}

	f, err := os.Create(path.Join(tag.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		errCh <- err
		return
	}

	if err := f.Chmod(0644); err != nil {
		errCh <- err
		return
	}

	if err := w.layout.RenderTag(f, tag); err != nil {
		errCh <- err
		return
	}

	doneCh <- true
}

func (w *Writer) writePost(post *pages.Post, doneCh chan<- bool, errCh chan<- error) {
	if err := os.MkdirAll(post.OutputPath, 0755); err != nil {
		errCh <- err
		return
	}

	if err := copyExclude(post.InputPath, post.OutputPath, ".md"); err != nil {
		errCh <- err
		return
	}

	f, err := os.Create(path.Join(post.OutputPath, "/index.html"))
	defer f.Close()
	if err != nil {
		errCh <- err
		return
	}

	if err := f.Chmod(0644); err != nil {
		errCh <- err
		return
	}

	if err := w.layout.RenderPost(f, post); err != nil {
		errCh <- err
		return
	}

	doneCh <- true
}
