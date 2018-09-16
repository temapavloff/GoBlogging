package layout

import (
	"GoBlogging/config"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	txt "text/template"
	"time"
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

	path := getTemplatePath(c, "")
	if _, err := os.Stat(path); err != nil {
		panic(fmt.Errorf("Template not found at %s", path))
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

// GetPostTpl - returns post template
func (l Layout) GetPostTpl() (*template.Template, error) {
	tpl, err := l.prepareLayout()
	if err != nil {
		return tpl, err
	}

	return tpl.Parse(l.postTemplate)
}

// GetIndexTpl - returns index template
func (l Layout) GetIndexTpl() (*template.Template, error) {
	tpl, err := l.prepareLayout()
	if err != nil {
		return tpl, err
	}

	return tpl.Parse(l.indexTemplate)
}

// GetTagTpl - return tag template
func (l Layout) GetTagTpl() (*template.Template, error) {
	tpl, err := l.prepareLayout()
	if err != nil {
		return tpl, err
	}

	return tpl.Parse(l.tagTemplate)
}

// GetRSSTpl - returns template for RSS-feed
func (l Layout) GetRSSTpl() (*txt.Template, error) {
	tpl := txt.New("layout").Funcs(txt.FuncMap{
		"formatDateRFC": func(date time.Time) string {
			return date.Format(time.RFC1123Z)
		},
	})
	return tpl.Parse(rssDefault)
}

func getTemplatePath(c *config.Config, name string) string {
	return c.GetAbsPath("/templates/" + c.Template + "/" + name)
}

func loadTemplate(c *config.Config, name, fallback string) string {
	path := getTemplatePath(c, name)
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Template file %s not found. Falling back to default template.\n", name)
		return fallback
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}
