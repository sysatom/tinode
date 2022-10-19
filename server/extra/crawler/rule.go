package crawler

import (
	"bytes"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Rule struct {
	Name    string
	Id      string
	Channel string
	When    string
	Mode    string
	Page    struct {
		URL  string
		List string
		Item map[string]string
	}
}

func (r Rule) Run() [][]byte {
	var result [][]byte

	doc, err := document(r.Page.URL)
	if err != nil {
		return result
	}

	doc.Find(r.Page.List).Each(func(i int, s *goquery.Selection) {
		// sort keys
		keys := make([]string, 0, len(r.Page.Item))
		for k := range r.Page.Item {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		txt := bytes.Buffer{}
		for _, k := range keys {
			f := ParseFun(s, r.Page.Item[k])
			v, err := f.Invoke()
			if err != nil {
				continue
			}
			v = strings.TrimSpace(v)
			v = strings.ReplaceAll(v, "\n", "")
			v = strings.ReplaceAll(v, "\r\n", "")
			if v == "" {
				continue
			}
			txt.WriteString(k)
			txt.WriteString(": ")
			txt.WriteString(v)
			txt.WriteString("\n")
		}
		if txt.Len() == 0 {
			return
		}
		result = append(result, txt.Bytes())
	})
	return result
}

type Result struct {
	Name    string
	Channel string
	Mode    string
	Result  [][]byte
}

func document(url string) (*goquery.Document, error) {
	res, err := http.Get(url) // #nosec
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	return goquery.NewDocumentFromReader(res.Body)
}
