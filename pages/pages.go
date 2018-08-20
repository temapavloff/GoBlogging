package pages

import (
	"GoBlogging/config"
	"sort"
)

// Pages - representation of whole blog
type Pages struct {
	config *config.Config
	Index  *Index
	Tags   *Tags
}

// New - creates new instance of Pages
func New(c *config.Config) *Pages {
	return &Pages{
		config: c,
		Index:  &Index{Title: c.BlogTitle},
		Tags:   &Tags{data: make(map[string]*Tag)},
	}
}

// Add - adds new post onto blog structure
func (p *Pages) Add(post *Post) {
	p.Index.addPost(post)
	p.Tags.updateTags(p.config.ServerPath, post)
}

// Index - representation of index page
type Index struct {
	Title string
	Posts []*Post
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

// Tag - representation of tag page
type Tag struct {
	Title string
	Count int
	Posts []*Post
	URL   string
}

// Tags - representation of all tags
type Tags struct {
	data map[string]*Tag
}

func (t *Tags) updateTags(serverPath string, p *Post) {
	for _, tagString := range p.StringTags {
		p.Tags = append(p.Tags, t.updateOneTag(serverPath, tagString, p))
	}
}

func (t *Tags) updateOneTag(serverPath, tagString string, p *Post) *Tag {
	if _, has := t.data[tagString]; !has {
		newTag := &Tag{Title: tagString, URL: serverPath + "/" + tagString}
		t.data[tagString] = newTag
	}

	tag := t.data[tagString]
	tag.Count++
	tag.Posts = append(tag.Posts, p)
	t.data[tagString] = tag

	return tag
}
