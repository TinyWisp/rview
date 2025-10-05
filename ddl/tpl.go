package ddl

import (
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
		def               *regexp.Regexp
		vfor              *regexp.Regexp
		whitespace        *regexp.Regexp
	}{
		openingTagBegin:   regexp.MustCompile(`^<([a-zA-Z0-9\-]+)`),
		openingTagEnd:     regexp.MustCompile(`^>`),
		closingTag:        regexp.MustCompile(`^</([a-zA-Z0-9\-]+)>`),
		selfClosingTagEnd: regexp.MustCompile(`^/>`),
		attrStart:         regexp.MustCompile(`^([a-zA-Z0-9\-_@:]+)=`),
		attrWithoutVal:    regexp.MustCompile(`^([a-zA-Z0-9\-]+)`),
		def:               regexp.MustCompile(`^([a-zA-Z0-9_\-]+)\((.*)\)$`),
		vfor:              regexp.MustCompile(`^\s*\(\s*([a-zA-Z\_][a-zA-Z0-9\_]*),\s*([a-zA-Z\_][a-zA-Z0-9\_]*)\s*\)\s+of\s+(.*?)\s*$`),
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
	Parent     *TplNode
	Children   []*TplNode
	Idx        int
	Events     map[string]*TplAttr
	Attrs      map[string]*TplAttr
	Directives map[string]string
	Text       string
	Exp        *Exp
	If         *TplAttr
	ElseIf     *TplAttr
	Else       *TplAttr
	For        *TplFor
	Def        *TplAttr
	Pos        int
}

type TplAttr struct {
	Pos int
	Exp *Exp
}

type TplFor struct {
	Pos      int
	Idx      string
	IdxPos   int
	Val      string
	ValPos   int
	Range    *Exp
	RangePos int
}

func (tn *TplNode) addAttr(pos int, key string, val string) error {
	// def
	if key == "def" {
		matches := tplPattern.def.FindStringSubmatch(val)
		if len(matches) == 0 {
			return NewDdlError("", pos, "tpl.invalidDefAttr")
		}
		exp := &Exp{
			Type:       ExpFunc,
			FuncName:   matches[1],
			FuncParams: []*Exp{},
		}
		if len(matches) == 3 && trim(matches[2]) != "" {
			str := trim(matches[2])
			params := strings.Split(str, ",")
			for _, param := range params {
				pexp, perr := ParseExp(param)
				if perr != nil {
					return perr
				}
				if pexp.Type != ExpVar {
					return NewDdlError("", pos, "tpl.invalidDefAttr")
				}
				exp.FuncParams = append(exp.FuncParams, pexp)
			}
		}
		tn.Def = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// v-if
	} else if key == "v-if" {
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.If != nil {
			return NewDdlError("", pos, "tpl.duplicateDirective")
		}
		if tn.ElseIf != nil || tn.Else != nil || tn.For != nil {
			return NewDdlError("", pos, "tpl.conflictedDirective")
		}
		tn.If = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// v-else-if
	} else if key == "v-else-if" {
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.ElseIf != nil {
			return NewDdlError("", pos, "tpl.duplicateDirective")
		}
		if tn.If != nil || tn.Else != nil || tn.For != nil {
			return NewDdlError("", pos, "tpl.conflictedDirective")
		}
		tn.ElseIf = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// v-else
	} else if key == "v-else" {
		if tn.Else != nil {
			return NewDdlError("", pos, "tpl.duplicateDirective")
		}
		if tn.If != nil || tn.ElseIf != nil || tn.For != nil {
			return NewDdlError("", pos, "tpl.conflictedDirective")
		}
		tn.Else = &TplAttr{
			Pos: pos,
			Exp: &Exp{
				Type: ExpNil,
			},
		}

		// v-for
	} else if key == "v-for" {
		if tn.For != nil {
			return NewDdlError("", pos, "tpl.duplicateDirective")
		}
		if tn.If != nil || tn.ElseIf != nil || tn.Else != nil {
			return NewDdlError("", pos, "tpl.conflictedDirective")
		}
		matches := tplPattern.vfor.FindStringSubmatch(val)
		if len(matches) == 0 {
			return NewDdlError("", pos, "tpl.invalidVforDirective")
		}
		rangeExp, err := ParseExp(matches[3])
		if err != nil {
			return err
		}
		imatches := tplPattern.vfor.FindStringSubmatchIndex(val)
		tn.For = &TplFor{
			Pos:      pos,
			Idx:      matches[1],
			Val:      matches[2],
			Range:    rangeExp,
			IdxPos:   imatches[2],
			ValPos:   imatches[4],
			RangePos: imatches[6],
		}

		// v-bind:var
	} else if strings.HasPrefix(key, "v-bind:") {
		vname := key[7:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Attrs == nil {
			tn.Attrs = make(map[string]*TplAttr)
		}
		if _, ok := tn.Attrs[vname]; ok {
			return NewDdlError("", pos, "tpl.duplicateAttribute")
		}
		tn.Attrs[vname] = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// :var
	} else if strings.HasPrefix(key, ":") {
		vname := key[1:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Attrs == nil {
			tn.Attrs = make(map[string]*TplAttr)
		}
		if _, ok := tn.Attrs[vname]; ok {
			return NewDdlError("", pos, "tpl.duplicateAttribute")
		}
		tn.Attrs[vname] = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// v-on:event
	} else if strings.HasPrefix(key, "v-on:") {
		event := key[5:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*TplAttr)
		}
		if _, ok := tn.Events[event]; ok {
			return NewDdlError("", pos, "tpl.duplicateEventHandler")
		}
		tn.Events[event] = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// @event
	} else if strings.HasPrefix(key, "@") {
		event := key[1:]
		exp, err := ParseExp(val)
		if err != nil {
			return err
		}
		if tn.Events == nil {
			tn.Events = make(map[string]*TplAttr)
		}
		if _, ok := tn.Events[event]; ok {
			return NewDdlError("", pos, "tpl.duplicateEventHandler")
		}
		tn.Events[event] = &TplAttr{
			Pos: pos,
			Exp: exp,
		}

		// ordinary attritube
	} else {
		if tn.Attrs == nil {
			tn.Attrs = make(map[string]*TplAttr)
		}
		if _, ok := tn.Attrs[key]; ok {
			return NewDdlError("", pos, "tpl.duplicateAttribute")
		}
		tn.Attrs[key] = &TplAttr{
			Pos: pos,
			Exp: &Exp{
				Type: ExpStr,
				Str:  val,
			},
		}
	}

	return nil
}

func parseTpl(tpl string) ([]*TplNode, error) {
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

		nodeIdx := 0
		if parentTagNode != nil {
			nodeIdx = len(parentTagNode.Children)
		}

		left := tpl[pos:]

		// <tag
		if matches := tplPattern.openingTagBegin.FindStringSubmatch(left); !isReadingTag && len(matches) > 0 {
			name := matches[1]
			tagNode := TplNode{
				Type:    TplNodeTag,
				TagName: name,
				Pos:     pos,
				Parent:  parentTagNode,
				Idx:     nodeIdx,
			}
			tagNodeStack = append(tagNodeStack, &tagNode)
			isReadingTag = true
			if parentTagNode != nil {
				if trim(text) != "" {
					textNode := TplNode{
						Type:   TplNodeText,
						Text:   trim(text),
						Pos:    pos - len(text),
						Parent: parentTagNode,
						Idx:    nodeIdx,
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
				return nodeArr, NewDdlError(tpl, pos, "tpl.missingOpeningTag")
			}
			tagName := matches[1]
			if tagName != parentTagNode.TagName {
				return nodeArr, NewDdlError(tpl, pos, "tpl.mismatchedTag")
			}
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type:   TplNodeText,
					Text:   trim(text),
					Pos:    pos - len(text),
					Parent: parentTagNode,
					Idx:    nodeIdx,
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
			beginPos := pos
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
					return nodeArr, NewDdlError(tpl, pos+klen+1, "tpl.mismatchedDoubleQuotationMark")
				} else {
					return nodeArr, NewDdlError(tpl, pos+klen+1, "tpl.mismatchedSingleQuotationMark")
				}
			}
			err := curTagNode.addAttr(beginPos, attrKey, attrVal)
			if err != nil {
				if tpe, ok := err.(*DdlError); ok {
					tpe.SetDdl(tpl)
				}
				return nodeArr, err
			}

			// attritube without value,  like the "enabled" attribute in "<comp enabled>".
		} else if matches := tplPattern.attrWithoutVal.FindStringSubmatch(left); isReadingTag && len(matches) > 0 {
			attrKey := matches[1]
			attrVal := ""
			err := curTagNode.addAttr(pos, attrKey, attrVal)
			if err != nil {
				if tpe, ok := err.(*DdlError); ok {
					tpe.SetDdl(tpl)
					tpe.SetPos(pos)
				}
				return nodeArr, err
			}
			pos += len(matches[0])

			// {{ ... }}
		} else if !isReadingTag && strings.HasPrefix(left, "{{") {
			if parentTagNode != nil && trim(text) != "" {
				textNode := TplNode{
					Type:   TplNodeText,
					Text:   trim(text),
					Pos:    pos,
					Parent: parentTagNode,
					Idx:    nodeIdx,
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
						Type:   TplNodeExp,
						Exp:    exp,
						Pos:    pos + 2,
						Parent: parentTagNode,
						Idx:    nodeIdx,
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
			return nodeArr, NewDdlError(tpl, pos, "tpl.unexpectedToken")
		}
	}

	if isReadingTag && curTagNode != nil {
		return nodeArr, NewDdlError(tpl, curTagNode.Pos, "tpl.incompleteTag")
	}

	if len(tagNodeStack) > 0 {
		return nodeArr, NewDdlError(tpl, tagNodeStack[len(tagNodeStack)-1].Pos, "tpl.missingClosingTag")
	}

	return nodeArr, nil
}
