package builder

import (
	"GoBlogging/pages"
	"fmt"
	"io/ioutil"
	"path"
)

// ReaderFunc - worker function type declaration
type ReaderFunc func(*Builder, <-chan string, chan<- *pages.Post)

// Reader - default worker function
func Reader(b *Builder, pagesCh <-chan string, postCh chan<- *pages.Post) {
	for page := range pagesCh {
		pageContent, _ := ioutil.ReadFile(path.Join(page, "./index.md"))
		relPath := getRelativePath(b.config.GetAbsPath(b.config.Input), page)
		post, err := pages.NewPost(
			string(pageContent),
			page,
			b.config.GetAbsPath(b.config.Output+relPath),
			b.config.ServerPath+relPath)
		if err != nil {
			fmt.Printf("Cannot create Post object: %s\n", err)
		}
		b.mutex.Lock()
		b.pages.Add(post)
		b.mutex.Unlock()
		postCh <- post
	}
}

func getRelativePath(rootPath, pagePath string) string {
	rootLen := len(rootPath)
	return pagePath[rootLen:]
}
