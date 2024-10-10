package template

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	tplPattern = struct {
		openingTagBegin   *regexp.Regexp
		openingTagEnd     *regexp.Regexp
		closingTag        *regexp.Regexp
		selfClosingTagEnd *regexp.Regexp
		attrStart         *regexp.Regexp
		attrWithoutVal    *regexp.Regexp
		vfor              *regexp.Regexp
		whitespace        *regexp.Regexp
	}{
		openingTagBegin:   regexp.MustCompile(`^<([a-zA-Z0-9\-]+)`),
		openingTagEnd:     regexp.MustCompile(`^>`),
		closingTag:        regexp.MustCompile(`^</([a-zA-Z0-9\-]+)>`),
		selfClosingTagEnd: regexp.MustCompile(`^/>`),
		attrStart:         regexp.MustCompile(`^([a-zA-Z0-9\-_@:]+)=`),
		attrWithoutVal:    regexp.MustCompile(`^([a-zA-Z0-9\-]+)`),
		vfor:              regexp.MustCompile(`^ *([a-zA-Z\_][a-zA-Z0-9\_]*), *([a-zA-Z\_][a-zA-Z0-9\_]*) *:= *range +([a-zA-Z\_][a-zA-Z0-9\_]*) *$`),
		whitespace:        regexp.MustCompile(`^\s+`),
	}
)

type TplNodeType int

const (
	TplNodeTag TplNodeType = iota
	TplNodeText
	TplNodeExp
)

type TplNode struct {
	Type       TplNodeType
	TagName    string
	Children   []*TplNode
	Binds      map[string]*TplExp
	Events     map[string]*TplExp
	Attrs      map[string]string
	Directives map[string]string
	Text       string
	Exp        *TplExp
	If         *TplExp
	ElseIf     *TplExp
	Else       *TplExp
	For        *TplFor
}

type TplFor struct {
	Item  string
	Idx   string
	Range TplExp
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
		tn.Else = &TplExp{
			Type:     TplExpVar,
			Variable: "",
		}

		// v-for
	} else if key == "v-for" {
		matches := tplPattern.vfor.FindStringSubmatch(val)
		rangeExp, err := ParseTplExp(matches[3])
		if err != nil {
			return err
		}
		tn.For = &TplFor{
			Item:  matches[2],
			Idx:   matches[1],
			Range: *rangeExp,
		}

		// v-bind:var
	} else if strings.HasPrefix(key, "v-bind:") {
		vname := key[7:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		if tn.Binds == nil {
			tn.Binds = make(map[string]*TplExp)
		}
		tn.Binds[vname] = exp

		// :var
	} else if strings.HasPrefix(key, ":") {
		vname := key[1:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		if tn.Binds == nil {
			tn.Binds = make(map[string]*TplExp)
		}
		tn.Binds[vname] = exp

		// v-on:event
	} else if strings.HasPrefix(key, "v-on:") {
		event := key[5:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*TplExp)
		}
		tn.Events[event] = exp

		// @event
	} else if strings.HasPrefix(key, "@") {
		event := key[1:]
		exp, err := ParseTplExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*TplExp)
		}
		tn.Events[event] = exp

		// ordinary attritube
	} else {
		if tn.Attrs == nil {
			tn.Attrs = make(map[string]string)
		}
		tn.Attrs[key] = val
	}

	return nil
}

func ParseTpl(tpl string) ([]*TplNode, error) {
	curTagNode := (*TplNode)(nil)
	parentTagNode := (*TplNode)(nil)
	tagNodeStack := make([]*TplNode, 0)
	nodeArr := make([]*TplNode, 0)
	isReadingTag := false
	text := ""
	pos := 0

	for {
		if pos >= len(tpl) {
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

		left := tpl[pos:]

		// <tag
		if matches := tplPattern.openingTagBegin.FindStringSubmatch(left); !isReadingTag && len(matches) > 0 {
			name := matches[1]
			tagNode := TplNode{
				Type:    TplNodeTag,
				TagName: name,
			}
			tagNodeStack = append(tagNodeStack, &tagNode)
			isReadingTag = true
			if parentTagNode != nil {
				if trim(text) != "" {
					textNode := TplNode{
						Type: TplNodeText,
						Text: trim(text),
					}
					text = ""
					parentTagNode.Children = append(parentTagNode.Children, &textNode)
				}
				parentTagNode.Children = append(parentTagNode.Children, &tagNode)
			} else {
				nodeArr = append(nodeArr, &tagNode)
			}
			pos += len(matches[0])

			// >
		} else if matches := tplPattern.openingTagEnd.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			isReadingTag = false
			pos += 1

			// />
		} else if matches := tplPattern.selfClosingTagEnd.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			isReadingTag = false
			curTagNode = nil
			pos += 2
			tagNodeStack = tagNodeStack[:slen-1]

			// </tag>
		} else if matches := tplPattern.closingTag.FindStringSubmatch(left); !isReadingTag && len(matches) > 0 {
			tagName := matches[1]
			if tagName != parentTagNode.TagName {
				msg := fmt.Sprintf("opening and ending tag mismatch: %s", tagName)
				return nodeArr, newParseTplError(tpl, pos, msg)
			}
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type: TplNodeText,
					Text: trim(text),
				}
				text = ""
				parentTagNode.Children = append(parentTagNode.Children, &textNode)
			}
			curTagNode = nil
			tagNodeStack = tagNodeStack[:slen-1]
			pos += len(matches[0])

			// key="val" or key='val'
		} else if matches := tplPattern.attrStart.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			attrKey := matches[1]
			attrVal := ""
			complete := false
			klen := len(attrKey)
			if left[klen+1] == '"' {
				for idx := klen + 2; idx < len(left); idx++ {
					if left[idx] == '"' && left[idx-1] != '\\' {
						attrVal = strings.ReplaceAll(string(left[klen+2:idx]), "\\\"", "\"")
						complete = true
						pos += idx + 1
						break
					}
				}
			} else if left[klen+1] == '\'' {
				for idx := klen + 2; idx < len(left); idx++ {
					if left[idx] == '\'' && left[idx-1] != '\\' {
						attrVal = strings.ReplaceAll(string(left[klen+2:idx]), "\\'", "'")
						complete = true
						pos += idx + 1
						break
					}
				}
			}
			if !complete {
				msg := fmt.Sprintf("uncomplete attribute:%s", attrKey)
				return nodeArr, newParseTplError(tpl, pos, msg)
			}
			curTagNode.setAttr(attrKey, attrVal)

			// attritube without value,  like the "enabled" attribute in "<comp enabled>".
		} else if matches := tplPattern.attrWithoutVal.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			attrKey := matches[1]
			attrVal := ""
			curTagNode.setAttr(attrKey, attrVal)
			pos += len(matches[0])

			// {{ ... }}
		} else if !isReadingTag && strings.HasPrefix(left, "{{") {
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type: TplNodeText,
					Text: trim(text),
				}
				text = ""
				parentTagNode.Children = append(parentTagNode.Children, &textNode)
			}

			inDoubleQuote := false
			inSingleQuote := false
			expStr := ""
			for idx := 2; idx < len(left)-2; idx++ {
				if left[idx] == '\'' && left[idx-1] != '\\' && !inDoubleQuote {
					inSingleQuote = !inSingleQuote
				} else if left[idx] == '"' && left[idx-1] != '\\' && !inSingleQuote {
					inDoubleQuote = !inDoubleQuote
				} else if !inSingleQuote && !inDoubleQuote && left[idx:idx+2] == "}}" {
					expStr = left[2:idx]
					exp, err := ParseTplExp(expStr)
					if err != nil {
						return nodeArr, err
					}
					expNode := TplNode{
						Type: TplNodeExp,
						Exp:  exp,
					}
					parentTagNode.Children = append(parentTagNode.Children, &expNode)
					pos += idx + 2
				}
			}

			// text
		} else if !isReadingTag && !strings.HasPrefix(left, "{{") {
			text += string(left[0])
			pos += 1

			// whitespace
		} else if matches := tplPattern.whitespace.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			pos += len(matches[0])

			// others
		} else {
			return nodeArr, fmt.Errorf("unexpected: %s", left)
		}
	}

	return nodeArr, nil
}
