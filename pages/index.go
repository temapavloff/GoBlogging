package pages

import (
	"GoBlogging/layout"
	"os"
	"path"
	"sort"
)

// Index - representation of index page
type Index struct {
	Title       string
	Posts       []*Post
	Output      string
	URL         string
	Description string
	AuthorName  string
	AuthorEmail string
	Lang        string
}

// AddPost - adds new post into pages
func (i *Index) addPost(p *Post) {
	i.Posts = append(i.Posts, p)
}

// Order - orders posts by creation time
func (i *Index) order() {
	sort.Slice(i.Posts, func(prev, next int) bool {
		return i.Posts[prev].Created.After(i.Posts[next].Created)
	})
}

func (i *Index) Write(l layout.Layout) error {
	f, err := os.Create(path.Join(i.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	tpl, err := l.GetIndexTpl()
	if err != nil {
		return err
	}

	return tpl.Execute(f, i)
}
