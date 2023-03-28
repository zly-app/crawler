package dom

import (
	"io"

	"github.com/antchfx/xmlquery"
)

type XmlDom struct {
	node *xmlquery.Node

	Type         xmlquery.NodeType
	Data         string
	Prefix       string
	NamespaceURI string
	Attr         []xmlquery.Attr
}

// 返回原始xmlquery.node
func (d *XmlDom) RawNode() *xmlquery.Node {
	return d.node
}

// 查找所有, 表达式错误会panic
func (d *XmlDom) XmlXpath(expr string) []*XmlDom {
	nodes := xmlquery.Find(d.node, expr)
	return makeXmlDom(nodes)
}

// 查找一个, 表达式错误会panic
func (d *XmlDom) XmlXpathOne(expr string) *XmlDom {
	node := xmlquery.FindOne(d.node, expr)
	return makeOneXmlDom(node)
}

// 获取属性
func (d *XmlDom) GetAttr(name string) string {
	return d.node.SelectAttr(name)
}

// 获取node内所有的文本值
func (d *XmlDom) InnerText() string {
	return d.node.InnerText()
}

/*
将node转为xml

	self 表示是否输入自己
*/
func (d *XmlDom) OutputXML(self bool) string {
	return d.node.OutputXML(self)
}

// 返回上级节点
func (d *XmlDom) Parent() *XmlDom {
	return makeOneXmlDom(d.node.Parent)
}

// 返回第一个子节点
func (d *XmlDom) FirstChild() *XmlDom {
	return makeOneXmlDom(d.node.FirstChild)
}

// 返回最后一个子节点
func (d *XmlDom) LastChild() *XmlDom {
	return makeOneXmlDom(d.node.LastChild)
}

// 返回上一个同级节点
func (d *XmlDom) PrevSibling() *XmlDom {
	return makeOneXmlDom(d.node.PrevSibling)
}

// 返回下一个同级节点
func (d *XmlDom) NextSibling() *XmlDom {
	return makeOneXmlDom(d.node.NextSibling)
}

// 获取所有子
func (d *XmlDom) Children() []*XmlDom {
	var a []*xmlquery.Node
	for nn := d.node.FirstChild; nn != nil; nn = nn.NextSibling {
		a = append(a, nn)
	}
	return makeXmlDom(a)
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
		dom[i] = makeOneXmlDom(node)
	}
	return dom
}

func makeOneXmlDom(node *xmlquery.Node) *XmlDom {
	if node == nil {
		return nil
	}
	return &XmlDom{
		node:         node,
		Type:         node.Type,
		Data:         node.Data,
		Prefix:       node.Prefix,
		NamespaceURI: node.NamespaceURI,
		Attr:         node.Attr,
	}
}
