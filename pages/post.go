package pages

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"

	"github.com/temapavloff/blackfriday"
)

// Post - representation of blog post
type Post struct {
	Title      string    `meta:"title"`
	Cover      string    `meta:"cover"`
	Created    time.Time `meta:"created"`
	Teaser     string    `meta:"teaser"`
	Tags       []*Tag
	Content    template.HTML
	InputPath  string
	OutputPath string
	URL        string

	StringTags []string `meta:"tags"`
}

// NewPost - creates new Post instance
func NewPost(fileContent, inputPath, outputPath, URL string) (*Post, error) {
	post := &Post{InputPath: inputPath, OutputPath: outputPath, URL: URL}
	err := parse(post, fileContent)
	if err != nil {
		return post, err
	}
	if post.Cover != "" {
		post.Cover = URL + "/" + post.Cover
	}
	return post, err
}

func parse(post *Post, fileData string) error {
	parts := strings.SplitN(strings.TrimSpace(fileData), "\n\n", 2)

	if len(parts) != 2 {
		return errors.New("Metadata and page content must be split by 2 newlines")
	}

	if err := parseMetadata(post, parts[0]); err != nil {
		return fmt.Errorf("Cannot parse post metadata: %s", err)
	}

	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	})
	var buf bytes.Buffer
	optList := []blackfriday.Option{
		blackfriday.WithRenderer(r),
		blackfriday.WithExtensions(blackfriday.CommonExtensions)}
	parser := blackfriday.New(optList...)
	ast := parser.Parse([]byte(parts[1]))
	r.RenderHeader(&buf, ast)
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Image && entering {
			dest := node.LinkData.Destination
			if dest[0] == '/' {
				dest = append([]byte(post.URL), dest...)
			}
			node.LinkData.Destination = dest
			node.Parent.Attributes.Add("class", "wide")
		}
		return blackfriday.GoToNext
	})
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(&buf, node, entering)
	})
	r.RenderFooter(&buf, ast)
	post.Content = template.HTML(buf.Bytes())

	return nil
}

func parseMetadata(post *Post, metaString string) error {
	metaMap := make(map[string]string)

	for _, line := range strings.Split(metaString, "\n") {
		keyVal := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(keyVal[0])

		if key != "" {
			metaMap[key] = strings.TrimSpace(keyVal[1])
		}
	}

	postRef := reflect.TypeOf(post).Elem()
	postVal := reflect.ValueOf(post).Elem()

	for i := 0; i < postVal.NumField(); i++ {
		fv := postVal.Field(i)
		ft := postRef.Field(i)
		metaName := ft.Tag.Get("meta")

		if _, has := metaMap[metaName]; fv.CanSet() && has {
			err := parseFieldByName(fv, metaName, metaMap[metaName])
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func parseFieldByName(field reflect.Value, name string, value string) error {
	if name == "tags" {
		tags := strings.Split(value, ",")
		for i, t := range tags {
			tags[i] = strings.TrimSpace(t)
		}
		field.Set(reflect.ValueOf(tags))
		return nil
	}

	if name == "created" {
		t, err := time.Parse("2006-01-02 15:04", value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(t))
		return nil
	}

	field.SetString(value)
	return nil
}
