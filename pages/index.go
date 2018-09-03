package pages

import (
	"html/template"
	"os"
	"path"
	"sort"
)

// Index - representation of index page
type Index struct {
	Title  string
	Posts  []*Post
	Output string
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

func (i *Index) Write(tpl *template.Template) error {
	f, err := os.Create(path.Join(i.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	return tpl.Execute(f, i)
}
