package template

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	openingTagBeginPattern   *regexp.Regexp = regexp.MustCompile(`^<([a-zA-Z0-9\-]+)`)
	openingTagEndPattern     *regexp.Regexp = regexp.MustCompile(`^>`)
	closingTagPattern        *regexp.Regexp = regexp.MustCompile(`^</([a-zA-Z0-9\-]+)>`)
	selfClosingTagEndPattern *regexp.Regexp = regexp.MustCompile(`^/>`)
	attrStartPattern         *regexp.Regexp = regexp.MustCompile(`^([a-zA-Z0-9\-@:])=`)
	forPattern               *regexp.Regexp = regexp.MustCompile(`^ *([a-zA-Z\_][a-zA-Z0-9\_]*), *([a-zA-Z\_][a-zA-Z0-9\_]*) *:= *range +([a-zA-Z\_][a-zA-Z0-9\_]*) *$`)
)

type TplNodeType int

const (
	TplNodeTag TplNodeType = iota
	TplNodeText
	TplNodeExp
)

type TplNode struct {
	Type            TplNodeType
	TagName         string
	Children        []*TplNode
	Binds           map[string]*TplExp
	Events          map[string]*TplExp
	Attrs           map[string]string
	Directives      map[string]string
	Text            string
	Exp             *TplExp
	If              *TplExp
	ElseIf          *TplExp
	Else            *TplExp
	ForItemVarName  string
	ForIdxVarName   string
	ForRangeVarName string
}

type ParseTplError struct {
	Tpl string
	Pos int
	Msg string
}

func (pte ParseTplError) Error() string {
	col := 0
	text := ""
	charNum := 0
	lines := strings.Split(pte.Tpl, "\n")
	for _, line := range lines {
		text += line + "\n"
		charNum += len(line) + 1
		if charNum >= pte.Pos {
			col = pte.Pos - (charNum - len(line) - 1)
			text += strings.Repeat(" ", col)
			break
		}
	}

	return fmt.Sprintf("%s\n%s", text, pte.Msg)
}

func newParseTplError(tpl string, pos int, msg string) ParseTplError {
	return ParseTplError{
		Tpl: tpl,
		Pos: pos,
		Msg: msg,
	}
}

func (tn *TplNode) setAttr(key string, val string) error {

	// v-if
	if key == "v-if" {
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		if tn.If != nil || tn.ElseIf != nil || tn.Else != nil {
			return fmt.Errorf("conflict")
		}
		tn.If = exp

		// v-else-if
	} else if key == "v-else-if" {
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.ElseIf = exp

		// v-else
	} else if key == "v-else" {
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.ElseIf = exp

		// v-for
	} else if key == "v-for" {
		matches := forPattern.FindStringSubmatch(val)
		tn.ForItemVarName = matches[2]
		tn.ForIdxVarName = matches[1]
		tn.ForRangeVarName = matches[3]

		// v-bind:var
	} else if strings.HasPrefix(key, "v-bind:") {
		vname := key[7:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.Binds[vname] = exp

		// :var
	} else if strings.HasPrefix(key, ":") {
		vname := key[1:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.Binds[vname] = exp

		// v-on:event
	} else if strings.HasPrefix(key, "v-on:") {
		event := key[5:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.Events[event] = exp

		// @event
	} else if strings.HasPrefix(key, "@") {
		event := key[1:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		tn.Events[event] = exp
	} else {
		tn.Attrs[key] = val
	}

	return nil
}

func ReadTpl(tpl string) ([]*TplNode, error) {
	byteArr := []byte(tpl)
	curTagNode := (*TplNode)(nil)
	parentTagNode := (*TplNode)(nil)
	tagNodeStack := make([]*TplNode, 0)
	nodeArr := make([]*TplNode, 0)
	isReadingTag := false
	text := ""
	pos := 0

	for {
		if pos >= len(byteArr) {
			break
		}

		slen := len(tagNodeStack)
		if isReadingTag {
			curTagNode = tagNodeStack[slen-1]
			if slen >= 2 {
				parentTagNode = tagNodeStack[slen-2]
			} else {
				parentTagNode = nil
			}
		} else {
			curTagNode = nil
			if slen > 0 {
				parentTagNode = tagNodeStack[slen-1]
			} else {
				parentTagNode = nil
			}
		}

		left := string(byteArr[pos:])

		// <tag
		if matches := openingTagBeginPattern.FindStringSubmatch(left); !isReadingTag && len(matches) > 0 {
			name := matches[1]
			tagNode := TplNode{
				Type:    TplNodeTag,
				TagName: name,
			}
			tagNodeStack = append(tagNodeStack, &tagNode)
			isReadingTag = true
			if parentTagNode != nil {
				textNode := TplNode{
					Type: TplNodeText,
					Text: text,
				}
				text = ""
				parentTagNode.Children = append(parentTagNode.Children, &textNode, &tagNode)
			} else {
				nodeArr = append(nodeArr, &tagNode)
			}
			pos += len(matches[0])

			// >
		} else if matches := openingTagEndPattern.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			isReadingTag = false
			pos += 1
			tagNodeStack = tagNodeStack[:slen-1]

			// />
		} else if matches := selfClosingTagEndPattern.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			isReadingTag = false
			curTagNode = nil
			pos += 2
			tagNodeStack = tagNodeStack[:slen-1]

			// </tag>
		} else if matches := closingTagPattern.FindStringSubmatch(left); !isReadingTag && len(matches) > 0 {
			tagName := matches[1]
			if tagName != curTagNode.TagName {
				msg := fmt.Sprintf("opening and ending tag mismatch: %s", tagName)
				return nodeArr, newParseTplError(tpl, pos, msg)
			}
			curTagNode = nil
			tagNodeStack = tagNodeStack[:slen-1]

			// key="val" or key='val'
		} else if matches := attrStartPattern.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			pos += len(matches[0])
			attrKey := matches[1]
			attrVal := ""
			complete := false
			if left[pos] == '"' {
				for idx := pos + 1; idx < len(left); idx++ {
					if left[idx] == '"' && left[idx-1] != '\\' {
						attrVal = strings.ReplaceAll(string(left[pos+1:idx]), "\\\"", "\"")
						complete = true
						break
					}
				}
			} else if left[pos] == '\'' {
				for idx := pos + 1; idx < len(left); idx++ {
					if left[idx] == '\'' && left[idx-1] != '\\' {
						attrVal = strings.ReplaceAll(string(left[pos+1:idx]), "\\'", "'")
						complete = true
						break
					}
				}
			}
			if !complete {
				msg := fmt.Sprintf("uncomplete attribute:%s", attrKey)
				return nodeArr, newParseTplError(tpl, pos, msg)
			}
			curTagNode.setAttr(attrKey, attrVal)

			// {{ ... }}
		} else if !isReadingTag && strings.HasPrefix(left, "{{") {
			if parentTagNode != nil {
				textNode := TplNode{
					Type: TplNodeText,
					Text: text,
				}
				text = ""
				parentTagNode.Children = append(parentTagNode.Children, &textNode)
			}

			inDoubleQuote := false
			inSingleQuote := false
			expStr := ""
			for idx := pos + 3; idx < len(left)-2; idx++ {
				if left[idx] == '\'' && left[idx-1] != '\\' && !inDoubleQuote {
					inSingleQuote = !inSingleQuote
				} else if left[idx] == '"' && left[idx-1] != '\\' && !inSingleQuote {
					inDoubleQuote = !inDoubleQuote
				} else if !inSingleQuote && !inDoubleQuote && left[idx:idx+2] == "}}" {
					expStr = left[pos+3 : idx]
					exp, err := ParseTplExp(expStr)
					if err != nil {
						return nodeArr, err
					}
					expNode := TplNode{
						Type: TplNodeExp,
						Exp:  exp,
					}
					parentTagNode.Children = append(parentTagNode.Children, &expNode)
					pos = idx + 2
				}
			}

			// text
		} else if !isReadingTag && !strings.HasPrefix(left, "{{") {
			text += string(left[0])
			pos += 1

			// others
		} else {
			return nodeArr, fmt.Errorf("unexpected: %s", left)
		}
	}

	return nodeArr, nil
}
