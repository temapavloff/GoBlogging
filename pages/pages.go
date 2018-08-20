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
		Index:  &Index{Title: c.BlogTitle, Output: c.GetAbsPath(c.Output)},
		Tags:   &Tags{data: make(map[string]*Tag)},
	}
}

// NodeType - type of page
type NodeType uint8

const (
	// IndexType - type of index page
	IndexType = 0
	// TagType - type of tag page
	TagType = 1
	// PostType - type of post page
	PostType = 2
)

// Node - special representation page to walk
type Node struct {
	Type  NodeType
	Index *Index
	Tag   *Tag
	Post  *Post
}

// PageWalker - type of callback for Pages.Walk
type PageWalker func(Node) error

// Walk - walks over all generated pages
func (p *Pages) Walk(walker PageWalker) error {
	p.Index.order()
	if err := walker(Node{Type: IndexType, Index: p.Index}); err != nil {
		return err
	}

	for _, post := range p.Index.Posts {
		if err := walker(Node{Type: PostType, Post: post}); err != nil {
			return err
		}
	}

	for _, tag := range p.Tags.data {
		tag.order()
		if err := walker(Node{Type: TagType, Tag: tag}); err != nil {
			return err
		}
	}

	return nil
}

// Len - reterns count of all pages
func (p *Pages) Len() int {
	return len(p.Index.Posts) + len(p.Tags.data) + 1
}

// Add - adds new post onto blog structure
func (p *Pages) Add(post *Post) {
	p.Index.addPost(post)
	p.Tags.updateTags(p.config, post)
}

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

// Tag - representation of tag page
type Tag struct {
	Title  string
	Count  int
	Posts  []*Post
	URL    string
	Output string
}

// Tags - representation of all tags
type Tags struct {
	data map[string]*Tag
}

func (t *Tag) order() {
	sort.Slice(t.Posts, func(prev, next int) bool {
		return t.Posts[prev].Created.After(t.Posts[next].Created)
	})
}

func (t *Tags) updateTags(c *config.Config, p *Post) {
	for _, tagString := range p.StringTags {
		p.Tags = append(p.Tags, t.updateOneTag(c, tagString, p))
	}
}

func (t *Tags) updateOneTag(c *config.Config, tagString string, p *Post) *Tag {
	if _, has := t.data[tagString]; !has {
		tagSlug := slug(tagString)
		newTag := &Tag{
			Title:  tagString,
			URL:    c.ServerPath + "/" + tagSlug,
			Output: c.GetAbsPath(c.Output + "/" + tagSlug),
		}
		t.data[tagString] = newTag
	}

	tag := t.data[tagString]
	tag.Count++
	tag.Posts = append(tag.Posts, p)
	t.data[tagString] = tag

	return tag
}
