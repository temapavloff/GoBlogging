package pages

import (
	"GoBlogging/config"
	"GoBlogging/layout"
)

// Page rendering iterface
type Page interface {
	Write(layout.Layout) error
}

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
		Index: &Index{
			Title:       c.BlogTitle,
			Output:      c.GetAbsPath(c.Output),
			URL:         c.Origin + c.ServerPath,
			Description: c.BlogDescription,
			AuthorName:  c.AuthorName,
			AuthorEmail: c.AuthorEmail,
			Lang:        c.Lang,
		},
		Tags: &Tags{data: make(map[string]*Tag)},
	}
}

// PageWalker - type of callback for Pages.Walk
type PageWalker func(Page) error

// Walk - walks over all generated pages
func (p *Pages) Walk(walker PageWalker) error {
	p.Index.order()
	if err := walker(p.Index); err != nil {
		return err
	}

	if err := walker(&RSS{p.Index}); err != nil {
		return err
	}

	for _, post := range p.Index.Posts {
		if err := walker(post); err != nil {
			return err
		}
	}

	for _, tag := range p.Tags.data {
		tag.order()
		if err := walker(tag); err != nil {
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
