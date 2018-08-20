package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"GoBlogging/config"
	"GoBlogging/layout"
	"GoBlogging/pages"
)

// Builder - main object for reading directory tree
type Builder struct {
	config  *config.Config
	workers int
	mutex   sync.Mutex
	pages   *pages.Pages
}

// New - creates new Reader instance
func New(c *config.Config) *Builder {
	return &Builder{
		c,
		runtime.NumCPU() * 2,
		sync.Mutex{},
		pages.New(c),
	}
}

// Run - starts build process
func (b *Builder) Read(worker ReaderFunc) {
	fmt.Printf("Reading directiry tree...\n")

	pagesCh := make(chan string)
	resultCh := make(chan *pages.Post)

	defer close(pagesCh)
	defer close(resultCh)

	for i := 0; i < b.workers; i++ {
		go worker(b, pagesCh, resultCh)
	}

	total := readTree(b.config.GetAbsPath(b.config.Input), pagesCh)

	cnt := 0
	for range resultCh {
		if cnt++; cnt == total {
			break
		}
	}

	fmt.Printf("Build blog structure, %d posts handled.\n", total)
}

func (b *Builder) Write() {
	fmt.Printf("Rendering...\n")

	l := layout.New()

	// l.RenderIndex(os.Stdout, b.pages.Index)
	for _, p := range b.pages.Index.Posts {
		l.RenderPost(os.Stdout, p)
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
