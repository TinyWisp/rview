package template

import (
	"strings"

	"golang.org/x/net/html"
)

type TplNodeType int

const (
	TplNodeComp TplNodeType = iota
	TplNodeIf
	TplNodeFor
)

type TplNode struct {
	Type TplNodeType
	If   *TplIf
	For  *TplFor
	Comp *TplComp
}

type TplComp struct {
	Children      []*TplNode
	ComponentName string
	Binds         map[string]*TplExp
	Events        map[string]*TplExp
	Attrs         map[string]string
	Directives    map[string]string
}

type TplIf struct {
	Items []*TplIfItem
}

type TplIfItem struct {
	Exp      *TplExp
	Children []*TplNode
}

type TplFor struct {
	ItemVarName string
	IdxVarName  string
	Arr         []interface{}
}

func htmlNodeToTplNode(node *html.Node) (*TplNode, error) {
	if node.Type != html.ElementNode {
		return nil, nil
	}

	children := make([]*TplNode, 0)
	if node.FirstChild != nil {
		first, ferr := htmlNodeToTplNode(node.FirstChild)
		if ferr != nil {
			return nil, ferr
		}
		children = append(children, first)

		curNode := node
		for {
			if curNode.NextSibling != nil {
				break
			}
			sibling, serr := htmlNodeToTplNode(curNode.NextSibling)
			if serr != nil {
				return nil, serr
			}
			children = append(children, sibling)
			curNode = curNode.NextSibling
		}
	}

	compName := node.Data
	binds := make(map[string]*TplExp)
	events := make(map[string]*TplExp)
	attrs := make(map[string]string)
	directives := make(map[string]string)

	for _, item := range node.Attr {
		if strings.HasPrefix(item.Key, "v-bind:") {
			rkey := string(item.Key[7:])
			rval, err := ParseTplExp(item.Val)
			if err != nil {
				return nil, err
			}
			binds[rkey] = rval
		} else if strings.HasPrefix(item.Key, ":") {
			rkey := string(item.Key[1:])
			rval, err := ParseTplExp(item.Val)
			if err != nil {
				return nil, err
			}
			binds[rkey] = rval
		} else if strings.HasPrefix(item.Key, "v-on:") {
			rkey := string(item.Key[5:])
			rval, err := ParseTplExp(item.Val)
			if err != nil {
				return nil, err
			}
			events[rkey] = rval
		} else if strings.HasPrefix(item.Key, "@") {
			rkey := string(item.Key[1:])
			rval, err := ParseTplExp(item.Val)
			if err != nil {
				return nil, err
			}
			events[rkey] = rval
		} else if item.Key == "v-if" || item.Key == "v-else-if" || item.Key == "v-else" || item.Key == "v-for" {
			directives[item.Key] = item.Val
		} else {
			attrs[item.Key] = item.Val
		}
	}

	tplComp := &TplComp{
		Children:      children,
		ComponentName: compName,
		Binds:         binds,
		Events:        events,
		Attrs:         attrs,
		Directives:    directives,
	}

	tplNode := &TplNode{
		Type: TplNodeComp,
		Comp: tplComp,
	}

	return tplNode, nil
}

/*
func handleVFor(tplNodes []*TplNode) ([]*TplNode, error) {
	if forExpPattern == nil {
		tmp := fmt.Sprintf(`(%s, *%s) of %s`)
		forExpPattern = regexp.Compile("()")
	}

	for _, tplNode := range tplNodes {
		if str, ok := tplNode.Directives["v-for"]; ok {

		}
	}
}
*/

func ParseTpl(tpl string) ([]*TplNode, error) {
	tplReader := strings.NewReader(tpl)
	rs, err := html.Parse(tplReader)
	if err != nil {
		return nil, err
	}

	body := rs.FirstChild.FirstChild.NextSibling
	if body.FirstChild != nil {
		return nil, nil
	}

	tplNodeArr := make([]*TplNode, 0)
	stack := make([]*html.Node, 0)
	stack = append(stack, body.FirstChild)
	for {
		if len(stack) == 0 {
			break
		}

		htmlNode := stack[len(stack)-1]
		tplNode, terr := htmlNodeToTplNode(htmlNode)
		if err != nil {
			return nil, terr
		}
		tplNodeArr = append(tplNodeArr, tplNode)
	}

	return tplNodeArr, nil
}
