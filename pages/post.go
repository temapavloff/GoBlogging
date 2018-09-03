package pages

import (
	"GoBlogging/helpers"
	"bufio"
	"bytes"
	"encoding/base64"
	"html/template"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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

	sourceContent []byte
}

// NewPost - creates new Post instance
func NewPost(inputPath, outputPath, URL string) (*Post, error) {
	post := &Post{InputPath: inputPath, OutputPath: outputPath, URL: URL}

	file, err := os.Open(inputPath + "/index.md")
	if err != nil {
		return post, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	metaMap, err := readMetadata(scanner)
	if err != nil {
		return post, err
	}

	if err := setMetadata(post, metaMap); err != nil {
		return post, err
	}

	src, err := readTail(scanner)
	if err != nil {
		return post, err
	}

	post.sourceContent = src
	if post.Cover != "" {
		post.Cover = URL + post.Cover
	}

	return post, err
}

func (p *Post) Write(tpl *template.Template) error {
	if err := os.MkdirAll(p.OutputPath, 0755); err != nil {
		return err
	}

	if err := helpers.CopyExclude(p.InputPath, p.OutputPath, ".md"); err != nil {
		return err
	}

	f, err := os.Create(path.Join(p.OutputPath, "/index.html"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}
	p.render()
	return tpl.Execute(f, p)
}

func (p *Post) render() {
	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags,
	})
	var buf bytes.Buffer
	optList := []blackfriday.Option{
		blackfriday.WithRenderer(r),
		blackfriday.WithExtensions(blackfriday.CommonExtensions)}
	parser := blackfriday.New(optList...)
	ast := parser.Parse(p.sourceContent)

	r.RenderHeader(&buf, ast)

	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Image && entering && node.Parent.Type == blackfriday.Paragraph {
			dest := string(node.LinkData.Destination)
			if dest[0] == '/' {
				base64img, w, h, err := writeImage(p.InputPath+dest, p.OutputPath+dest)
				if err != nil {
					// Just skip for now
					return blackfriday.GoToNext
				}
				node.Attributes.Add("class", "js-img blur")
				node.Attributes.Add("data-src", p.URL+dest)
				node.Attributes.Add("width", strconv.Itoa(w))
				node.Attributes.Add("height", strconv.Itoa(h))
				node.LinkData.Destination = []byte(base64img)
			}

			node.Parent.Attributes.Add("class", "wide")
		}
		return blackfriday.GoToNext
	})

	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(&buf, node, entering)
	})

	r.RenderFooter(&buf, ast)

	p.Content = template.HTML(buf.Bytes())
}

func writeImage(src, dst string) (string, int, int, error) {
	img, err := imaging.Open(src)
	if err != nil {
		return "", 0, 0, nil
	}

	srcW := img.Bounds().Size().X
	tagretW := 1200
	if tagretW > srcW {
		tagretW = srcW
	}

	destImg := imaging.Resize(img, tagretW, 0, imaging.NearestNeighbor)
	targetH := destImg.Bounds().Size().Y

	if err := imaging.Save(destImg, dst); err != nil {
		return "", 0, 0, nil
	}

	base64img := imaging.Resize(img, 100, 0, imaging.Linear)
	var buf bytes.Buffer
	imaging.Encode(&buf, base64img, imaging.JPEG)
	str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return "data:image/jpeg;base64," + str, tagretW, targetH, nil
}

func readMetadata(scanner *bufio.Scanner) (map[string]string, error) {
	result := make(map[string]string)
	var emptyLines uint8
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			emptyLines++
			if emptyLines == 2 {
				break
			}
			continue
		}
		if line != "" && emptyLines < 0 {
			emptyLines--
		}
		parts := strings.SplitN(line, ":", 2)
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func readTail(scanner *bufio.Scanner) ([]byte, error) {
	var buf []byte

	for scanner.Scan() {
		buf = append(buf, scanner.Bytes()...)
		buf = append(buf, []byte("\n")...)
	}

	if err := scanner.Err(); err != nil {
		return buf, err
	}

	return buf, nil
}

func setMetadata(p *Post, metaMap map[string]string) error {
	postRef := reflect.TypeOf(p).Elem()
	postVal := reflect.ValueOf(p).Elem()

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
