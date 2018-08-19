package pages

import (
	"sort"
)

// Index - representation of index page
type Index struct {
	Title string
	Posts []*Post
}

// NewIndex - creates new index pages instance
func NewIndex(title string) *Index {
	return &Index{Title: title}
}

// AddPost - adds new post into pages
func (i *Index) AddPost(p *Post) {
	i.Posts = append(i.Posts, p)
}

// Order - orders posts by creation time
func (i *Index) Order() {
	sort.Slice(i.Posts, func(prev, next int) bool {
		return i.Posts[prev].Created.After(i.Posts[next].Created)
	})
}

// Tag - representation of tag page
type Tag struct {
	Title string
	Count int
	Posts []*Post
	URL   string
}

// Tags - representation of all tags
type Tags struct {
	data map[string]Tag
}

// NewTags - creates new tags page instance
func NewTags() *Tags {
	return &Tags{data: make(map[string]Tag)}
}

// UpdateTags - updates tags data from given post
func (t *Tags) UpdateTags(serverPath string, p *Post) {
	for _, tagString := range p.Tags {
		t.updateOneTag(serverPath, tagString, p)
	}
}

// All - returns map of al tags
func (t *Tags) All() map[string]Tag {
	return t.data
}

func (t *Tags) updateOneTag(serverPath, tagString string, p *Post) {
	if _, has := t.data[tagString]; !has {
		t.data[tagString] = Tag{Title: tagString, URL: serverPath + "/" + tagString}
	}

	tag := t.data[tagString]
	tag.Count++
	tag.Posts = append(tag.Posts, p)
	t.data[tagString] = tag
}
