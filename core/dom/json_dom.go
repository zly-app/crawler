package dom

import (
	"io"

	"github.com/antchfx/jsonquery"
)

type JsonDom struct {
	node *jsonquery.Node

	Type jsonquery.NodeType
	Data string
}

// 返回原始jsonquery.node
func (j *JsonDom) RawNode() *jsonquery.Node {
	return j.node
}

// xpath查找所有, 表达式错误会panic
func (j *JsonDom) Xpath(expr string) []*JsonDom {
	nodes := jsonquery.Find(j.node, expr)
	return makeJsonDom(nodes)
}

// xpath查找所有, 表达式错误会panic
func (j *JsonDom) XpathOne(expr string) *JsonDom {
	node := jsonquery.FindOne(j.node, expr)
	return makeOneJsonDom(node)
}

// 查找具有指定name的第一个子
func (j *JsonDom) FindOneChild(name string) *JsonDom {
	node := j.node.SelectElement(name)
	return makeOneJsonDom(node)
}

// 获取node内所有的文本值
func (j *JsonDom) InnerText() string {
	return j.node.InnerText()
}

// 将node转为xml
func (j *JsonDom) OutputXML() string {
	return j.node.OutputXML()
}

// 返回上级节点
func (j *JsonDom) Parent() *JsonDom {
	return makeOneJsonDom(j.node.Parent)
}

// 返回第一个子节点
func (j *JsonDom) FirstChild() *JsonDom {
	return makeOneJsonDom(j.node.FirstChild)
}

// 返回最后一个子节点
func (j *JsonDom) LastChild() *JsonDom {
	return makeOneJsonDom(j.node.LastChild)
}

// 返回上一个同级节点
func (j *JsonDom) PrevSibling() *JsonDom {
	return makeOneJsonDom(j.node.PrevSibling)
}

// 返回下一个同级节点
func (j *JsonDom) NextSibling() *JsonDom {
	return makeOneJsonDom(j.node.NextSibling)
}

// 获取所有子
func (j *JsonDom) Children() []*JsonDom {
	nodes := j.node.ChildNodes()
	return makeJsonDom(nodes)
}

func NewJsonDom(r io.Reader) (*JsonDom, error) {
	node, err := jsonquery.Parse(r)
	if err != nil {
		return nil, err
	}
	return makeOneJsonDom(node), nil
}

func makeJsonDom(nodes []*jsonquery.Node) []*JsonDom {
	dom := make([]*JsonDom, len(nodes))
	for i, node := range nodes {
		dom[i] = makeOneJsonDom(node)
	}
	return dom
}

func makeOneJsonDom(node *jsonquery.Node) *JsonDom {
	if node == nil {
		return nil
	}
	return &JsonDom{
		node: node,
		Type: node.Type,
		Data: node.Data,
	}
}
