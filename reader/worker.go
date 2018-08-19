package reader

import (
	"GoBlogging/pages"
	"fmt"
	"io/ioutil"
	"path"
)

// WorkerFunc - worker function type declaration
type WorkerFunc func(*Reader, <-chan string, chan<- *pages.Post)

// Worker - default worker function
func Worker(r *Reader, pagesCh <-chan string, postCh chan<- *pages.Post) {
	for page := range pagesCh {
		pageContent, _ := ioutil.ReadFile(path.Join(page, "./index.md"))
		relPath := getRelativePath(r.config.GetAbsPath(r.config.Input), page)
		post, err := pages.NewPost(
			string(pageContent),
			page,
			r.config.GetAbsPath(r.config.Output+relPath),
			r.config.ServerPath+relPath)
		if err != nil {
			fmt.Printf("Cannot create Post object: %s\n", err)
		}
		r.mutex.Lock()
		r.Index.AddPost(post)
		r.Tags.UpdateTags(r.config.ServerPath, post)
		r.mutex.Unlock()
		postCh <- post
	}
}

func getRelativePath(rootPath, pagePath string) string {
	rootLen := len(rootPath)
	return pagePath[rootLen:]
}
