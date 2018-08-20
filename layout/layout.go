package layout

import (
	"GoBlogging/pages"
	"html/template"
	"io"
)

// Layout - blog template representation
type Layout struct {
	outerTemplate string
	postTemplate  string
	indexTemplate string
	tagTemplate   string
}

// New - creates new Layout instance
func New() Layout {

	return Layout{
		outerTemplate: outerDefault,
		postTemplate:  postDefault,
		indexTemplate: indexDefault,
		tagTemplate:   tagDefault,
	}
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
