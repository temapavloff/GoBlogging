package pages

import (
	"GoBlogging/layout"
	"os"
	"path"
)

// RSS - representation of index page
type RSS struct {
	*Index
}

func (r *RSS) Write(l layout.Layout) error {
	if err := os.MkdirAll(path.Join(r.Output, "/feeds"), 0755); err != nil {
		return err
	}

	f, err := os.Create(path.Join(r.Output, "/feeds", "/rss.xml"))
	defer f.Close()
	if err != nil {
		return err
	}

	if err := f.Chmod(0644); err != nil {
		return err
	}

	tpl, err := l.GetRSSTpl()
	if err != nil {
		return err
	}

	return tpl.Execute(f, r)
}
