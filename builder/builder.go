package builder

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

	var wg sync.WaitGroup
	pagesCh := make(chan string)

	defer close(pagesCh)

	for i := 0; i < b.workers; i++ {
		go worker(b, pagesCh, &wg)
	}

	total := readTree(b.config.GetAbsPath(b.config.Input), pagesCh, &wg)
	wg.Wait()

	fmt.Printf("Build blog structure, %d posts handled.\n", total)
}

func (b *Builder) Write(w *Writer) {
	fmt.Printf("Rendering...\n")

	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	nodeCh := make(chan pages.Node)

	defer close(nodeCh)

	if err := w.Prepare(); err != nil {
		panic(err)
	}
	for i := 0; i < b.workers; i++ {
		go w.Write(nodeCh, errCh, &wg)
	}
	b.pages.Walk(func(n pages.Node) error {
		wg.Add(1)
		nodeCh <- n
		return nil
	})

	go func() {
		wg.Wait()
		close(errCh)
		fmt.Printf("%d pages has been written.\n", b.pages.Len())
	}()

	select {
	case err, ok := <-errCh:
		if !ok {
			break
		}
		panic(err)
	}
}

func readTree(dir string, pages chan<- string, wg *sync.WaitGroup) int {
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

		wg.Add(1)
		total++
		pages <- curPath

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("Cannot read directory tree: %s", err))
	}
	return total
}
