package ddl

type DDLDef struct {
	TplMap      map[string]*TplNode
	CssClassMap CSSClassMap
}

func ParseDdl(ddl string) (DDLDef, error) {
	def := DDLDef{
		TplMap:      map[string]*TplNode{},
		CssClassMap: CSSClassMap{},
	}

	tplNodes, err := parseTpl(ddl)
	if err != nil {
		return def, err
	}
	for _, tn := range tplNodes {
		// template
		if tn.TagName == "template" {
			tplName := "main"
			if tn.Def != nil {
				tplName = tn.Def.Exp.FuncName
			}
			def.TplMap[tplName] = tn

			// style
		} else if tn.TagName == "style" {
			if len(tn.Children) > 1 {
				return def, NewDdlError(ddl, tn.Pos, "ddl.invalidStyleSection")
			}
			if len(tn.Children) == 1 && tn.Children[0].Type != TplNodeText {
				return def, NewDdlError(ddl, tn.Pos, "ddl.invalidStyleSection")
			}
			classMap, cerr := parseCss(tn.Children[0].Text)
			if cerr != nil {
				return def, cerr
			}
			for key, val := range classMap {
				def.CssClassMap[key] = val
			}
		}
	}

	return def, nil
}
