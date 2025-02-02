package template

import (
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	Binds      map[string]*Exp
	Events     map[string]*Exp
	Attrs      map[string]string
	Directives map[string]string
	Text       string
	Exp        *Exp
	If         *Exp
	ElseIf     *Exp
	Else       *Exp
	For        *TplFor
	Pos        int
}

type TplFor struct {
	Item  string
	Idx   string
	Range Exp
}

func (tn *TplNode) setAttr(key string, val string) error {
	// v-if
	if key == "v-if" {
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.If != nil {
			return NewTplParseError("", "tpl.duplicateDirective", 0)
		}
		if tn.ElseIf != nil || tn.Else != nil || tn.For != nil {
			return NewTplParseError("", "tpl.conflictedDirective", 0)
		}
		tn.If = exp

		// v-else-if
	} else if key == "v-else-if" {
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.ElseIf != nil {
			return NewTplParseError("", "tpl.duplicateDirective", 0)
		}
		if tn.If != nil || tn.Else != nil || tn.For != nil {
			return NewTplParseError("", "tpl.conflictedDirective", 0)
		}
		tn.ElseIf = exp

		// v-else
	} else if key == "v-else" {
		if tn.Else != nil {
			return NewTplParseError("", "tpl.duplicateDirective", 0)
		}
		if tn.If != nil || tn.ElseIf != nil || tn.For != nil {
			return NewTplParseError("", "tpl.conflictedDirective", 0)
		}
		tn.Else = &Exp{
			Type:     ExpVar,
			Variable: "",
		}

		// v-for
	} else if key == "v-for" {
		if tn.For != nil {
			return NewTplParseError("", "tpl.duplicateDirective", 0)
		}
		if tn.If != nil || tn.ElseIf != nil || tn.Else != nil {
			return NewTplParseError("", "tpl.conflictedDirective", 0)
		}
		matches := tplPattern.vfor.FindStringSubmatch(val)
		if len(matches) == 0 {
			return NewTplParseError("", "tpl.invalidForDirective", 0)
		}
		rangeExp, err := ParseExp(matches[3])
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
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Binds == nil {
			tn.Binds = make(map[string]*Exp)
		}
		if _, ok := tn.Binds[vname]; ok {
			return NewTplParseError("", "tpl.duplicateAttribute", 0)
		}
		if tn.Attrs != nil {
			if _, ok := tn.Attrs[vname]; ok {
				return NewTplParseError("", "tpl.duplicateAttribute", 0)
			}
		}
		tn.Binds[vname] = exp

		// :var
	} else if strings.HasPrefix(key, ":") {
		vname := key[1:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Binds == nil {
			tn.Binds = make(map[string]*Exp)
		}
		if _, ok := tn.Binds[vname]; ok {
			return NewTplParseError("", "tpl.duplicateAttribute", 0)
		}
		if tn.Attrs != nil {
			if _, ok := tn.Attrs[vname]; ok {
				return NewTplParseError("", "tpl.duplicateAttribute", 0)
			}
		}
		tn.Binds[vname] = exp

		// v-on:event
	} else if strings.HasPrefix(key, "v-on:") {
		event := key[5:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*Exp)
		}
		if _, ok := tn.Events[event]; ok {
			return NewTplParseError("", "tpl.duplicateEventHandler", 0)
		}
		tn.Events[event] = exp

		// @event
	} else if strings.HasPrefix(key, "@") {
		event := key[1:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*Exp)
		}
		if _, ok := tn.Events[event]; ok {
			return NewTplParseError("", "tpl.duplicateEventHandler", 0)
		}
		tn.Events[event] = exp

		// ordinary attritube
	} else {
		if tn.Attrs == nil {
			tn.Attrs = make(map[string]string)
		}
		spew.Dump(tn.Attrs)
		if _, ok := tn.Attrs[key]; ok {
			return NewTplParseError("", "tpl.duplicateAttribute", 0)
		}
		if tn.Binds != nil {
			if _, ok := tn.Binds[key]; ok {
				return NewTplParseError("", "tpl.duplicateAttribute", 0)
			}
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
				Pos:     pos,
			}
			tagNodeStack = append(tagNodeStack, &tagNode)
			isReadingTag = true
			if parentTagNode != nil {
				if trim(text) != "" {
					textNode := TplNode{
						Type: TplNodeText,
						Text: trim(text),
						Pos:  pos - len(text),
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
			if parentTagNode == nil {
				return nodeArr, NewTplParseError(tpl, "tpl.missingOpeningTag", pos)
			}
			tagName := matches[1]
			if tagName != parentTagNode.TagName {
				return nodeArr, NewTplParseError(tpl, "tpl.mismatchedTag", pos)
			}
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type: TplNodeText,
					Text: trim(text),
					Pos:  pos - len(text),
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
				if left[klen+1] == '"' {
					return nodeArr, NewTplParseError(tpl, "tpl.mismatchedDoubleQuotationMark", pos)
				} else {
					return nodeArr, NewTplParseError(tpl, "tpl.mismatchedSingleQuotationMark", pos)
				}
			}
			err := curTagNode.setAttr(attrKey, attrVal)
			if err != nil {
				if tpe, ok := err.(*TplParseError); ok {
					tpe.SetTpl(tpl)
					if tpe.IsExpError() {
						tpe.AddOffset(pos + klen + 2)
					} else {
						tpe.SetPos(pos)
					}
				}
				return nodeArr, err
			}

			// attritube without value,  like the "enabled" attribute in "<comp enabled>".
		} else if matches := tplPattern.attrWithoutVal.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			attrKey := matches[1]
			attrVal := ""
			err := curTagNode.setAttr(attrKey, attrVal)
			if err != nil {
				if tpe, ok := err.(*TplParseError); ok {
					tpe.SetTpl(tpl)
					tpe.SetPos(pos)
				}
				return nodeArr, err
			}
			pos += len(matches[0])

			// {{ ... }}
		} else if !isReadingTag && strings.HasPrefix(left, "{{") {
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type: TplNodeText,
					Text: trim(text),
					Pos:  pos,
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
					exp, err := ParseExp(expStr)
					if err != nil {
						return nodeArr, err
					}
					expNode := TplNode{
						Type: TplNodeExp,
						Exp:  exp,
						Pos:  pos + 2,
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
			return nodeArr, NewTplParseError(tpl, "tpl.unexpectedToken", pos)
		}
	}

	if isReadingTag && curTagNode != nil {
		return nodeArr, NewTplParseError(tpl, "tpl.incompleteTag", curTagNode.Pos)
	}

	if len(tagNodeStack) > 0 {
		return nodeArr, NewTplParseError(tpl, "tpl.missingClosingTag", tagNodeStack[len(tagNodeStack)-1].Pos)
	}

	return nodeArr, nil
}
