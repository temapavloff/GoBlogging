package builder

import (
	"GoBlogging/pages"
	"fmt"
	"sync"
)

// ReaderFunc - worker function type declaration
type ReaderFunc func(*Builder, <-chan string, *sync.WaitGroup)

// Reader - default worker function
func Reader(b *Builder, pagesCh <-chan string, wg *sync.WaitGroup) {
	for page := range pagesCh {
		relPath := getRelativePath(b.config.GetAbsPath(b.config.Input), page)
		post, err := pages.NewPost(
			page,
			b.config.GetAbsPath(b.config.Output+relPath),
			b.config.ServerPath+relPath)
		if err != nil {
			fmt.Printf("Cannot create Post object: %s\n", err)
		}
		b.mutex.Lock()
		b.pages.Add(post)
		b.mutex.Unlock()
		wg.Done()
	}
}

func getRelativePath(rootPath, pagePath string) string {
	rootLen := len(rootPath)
	return pagePath[rootLen:]
}
