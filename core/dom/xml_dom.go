package dom

import (
	"io"

	"github.com/antchfx/xmlquery"
)

type XmlDom struct {
	*xmlquery.Node
}

// 查找所有, 表达式错误会panic
func (d *XmlDom) XmlXpath(expr string) []*XmlDom {
	nodes := xmlquery.Find(d.Node, expr)
	return makeXmlDom(nodes)
}

// 查找一个, 表达式错误会panic
func (d *XmlDom) XmlXpathOne(expr string) *XmlDom {
	node := xmlquery.FindOne(d.Node, expr)
	return makeOneXmlDom(node)
}

func NewXmlDom(r io.Reader) (*XmlDom, error) {
	node, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}
	return makeOneXmlDom(node), nil
}

func makeXmlDom(nodes []*xmlquery.Node) []*XmlDom {
	dom := make([]*XmlDom, len(nodes))
	for i, node := range nodes {
		dom[i] = &XmlDom{node}
	}
	return dom
}

func makeOneXmlDom(node *xmlquery.Node) *XmlDom {
	return &XmlDom{node}
}
