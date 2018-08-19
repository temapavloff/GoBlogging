package reader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"GoBlogging/config"
	"GoBlogging/pages"
)

// Reader - main object for reading directory tree
type Reader struct {
	config      *config.Config
	workers     int
	startWorker WorkerFunc
	mutex       sync.Mutex
	Index       *pages.Index
	Tags        *pages.Tags
}

// New - creates new Reader instance
func New(c *config.Config, startWorker WorkerFunc) *Reader {
	return &Reader{
		c,
		runtime.NumCPU() * 2,
		startWorker,
		sync.Mutex{},
		pages.NewIndex(c.BlogTitle),
		pages.NewTags(),
	}
}

// Run - starts build process
func (r *Reader) Run() {
	pagesCh := make(chan string)
	resultCh := make(chan *pages.Post)

	defer close(pagesCh)
	defer close(resultCh)

	for i := 0; i < r.workers; i++ {
		go r.startWorker(r, pagesCh, resultCh)
	}

	total := readTree(r.config.GetAbsPath(r.config.Input), pagesCh)

	cnt := 0
	for range resultCh {
		if cnt++; cnt == total {
			break
		}
	}
}

func readTree(dir string, pages chan<- string) int {
	total := 0
	err := filepath.Walk(dir, func(curPath string, curInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}

		if !curInfo.IsDir() {
			return nil
		}

		pagePath := path.Join(curPath, "./index.md")

		if _, err := os.Stat(pagePath); os.IsNotExist(err) {
			return nil
		}

		total++
		pages <- curPath

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("Cannot read directory tree: %s", err))
	}

	return total
}
