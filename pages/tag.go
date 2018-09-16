package pages

import (
	"GoBlogging/config"
	"GoBlogging/layout"
	"os"
	"path"
	"sort"
)

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

func (t *Tag) Write(l layout.Layout) error {
	if err := os.MkdirAll(t.Output, 0755); err != nil {
		return err
	}

	f, err := os.Create(path.Join(t.Output, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	tpl, err := l.GetTagTpl()
	if err != nil {
		return err
	}

	return tpl.Execute(f, t)
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
			URL:    c.Origin + c.ServerPath + "/tags/" + tagSlug,
			Output: c.GetAbsPath(c.Output + "/tags/" + tagSlug),
		}
		t.data[tagString] = newTag
	}

	tag := t.data[tagString]
	tag.Count++
	tag.Posts = append(tag.Posts, p)
	t.data[tagString] = tag

	return tag
}
