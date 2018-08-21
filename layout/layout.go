package layout

import (
	"GoBlogging/config"
	"GoBlogging/pages"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

const outerName = "layout.html"
const postName = "post.html"
const indexName = "index.html"
const tagName = "tag.html"

// Layout - blog template representation
type Layout struct {
	config        *config.Config
	outerTemplate string
	postTemplate  string
	indexTemplate string
	tagTemplate   string
}

// New - creates new Layout instance
func New(c *config.Config) Layout {

	if c.Template == "" {
		return Layout{
			config:        c,
			outerTemplate: outerDefault,
			postTemplate:  postDefault,
			indexTemplate: indexDefault,
			tagTemplate:   tagDefault,
		}
	}

	if _, err := os.Stat(getTemplatePath(c, "")); err != nil {
		panic(err)
	}

	return Layout{
		config:        c,
		outerTemplate: loadTemplate(c, outerName, outerDefault),
		postTemplate:  loadTemplate(c, postName, postDefault),
		indexTemplate: loadTemplate(c, indexName, indexDefault),
		tagTemplate:   loadTemplate(c, tagName, tagDefault),
	}
}

// GetAssetsPath - returns assets path for selected
func (l Layout) GetAssetsPath() string {
	if l.config.Template == "" {
		return ""
	}

	return l.config.GetAbsPath("/templates/" + l.config.Template + "/assets")
}

func (l Layout) prepareLayout() (*template.Template, error) {
	tpl := template.New("layout")
	if _, err := tpl.Parse(l.outerTemplate); err != nil {
		return tpl, err
	}

	return tpl, nil
}

// RenderPost - renders single page layout
func (l Layout) RenderPost(writer io.Writer, post *pages.Post) error {
	tpl, err := l.prepareLayout()
	if err != nil {
		return err
	}
	if _, err := tpl.New("content").Parse(l.postTemplate); err != nil {
		return err
	}

	return tpl.Execute(writer, post)
}

// RenderIndex - renders main page
func (l Layout) RenderIndex(writer io.Writer, index *pages.Index) error {
	tpl, err := l.prepareLayout()
	if err != nil {
		return err
	}
	if _, err := tpl.New("content").Parse(l.indexTemplate); err != nil {
		return err
	}

	return tpl.Execute(writer, index)
}

// RenderTag - renders one tag page
func (l Layout) RenderTag(writer io.Writer, tag *pages.Tag) error {
	tpl, err := l.prepareLayout()
	if err != nil {
		return err
	}
	if _, err := tpl.New("content").Parse(l.tagTemplate); err != nil {
		return err
	}

	return tpl.Execute(writer, tag)
}

func getTemplatePath(c *config.Config, name string) string {
	return c.GetAbsPath("/templates/" + c.Template + "/" + name)
}

func loadTemplate(c *config.Config, name, fallback string) string {
	path := getTemplatePath(c, name)
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Template file %s not found. Falling back to default template.", name)
		return fallback
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}
