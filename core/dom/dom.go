package dom

import (
	"io"

	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Dom struct {
	*html.Node
}

// xpath查找所有, 表达式错误会panic
func (d *Dom) Xpath(expr string) []*Dom {
	nodes := htmlquery.Find(d.Node, expr)
	return makeDom(nodes)
}

// xpath查找一个, 表达式错误会panic
func (d *Dom) XpathOne(expr string) *Dom {
	node := htmlquery.FindOne(d.Node, expr)
	return makeOneDom(node)
}

// css查找所有, 表达式错误会panic
func (d *Dom) Css(expr string) []*Dom {
	sel := getCssQuery(expr)
	nodes := cascadia.QueryAll(d.Node, sel)
	return makeDom(nodes)
}

// css查找一个, 表达式错误会panic
func (d *Dom) CssOne(expr string) *Dom {
	sel := getCssQuery(expr)
	node := cascadia.Query(d.Node, sel)
	return makeOneDom(node)
}

// 获取属性
func (d *Dom) GetAttr(name string) string {
	return htmlquery.SelectAttr(d.Node, name)
}

// 获取node内所有的文本值
func (d *Dom) InnerText() string {
	return htmlquery.InnerText(d.Node)
}

/*将node转为html
  self 表示是否输入自己
*/
func (d *Dom) HTML(self bool) string {
	return htmlquery.OutputHTML(d.Node, self)
}

func NewDom(r io.Reader) (*Dom, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return makeOneDom(node), nil
}

func makeDom(nodes []*html.Node) []*Dom {
	dom := make([]*Dom, len(nodes))
	for i, node := range nodes {
		dom[i] = &Dom{node}
	}
	return dom
}

func makeOneDom(node *html.Node) *Dom {
	return &Dom{node}
}
