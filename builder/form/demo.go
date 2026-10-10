package form

func Demo() Form {
	choice := []Option{{Value: "a", Label: "选项A"}, {Value: "b", Label: "选项B"}}
	return Form{
		ID:    "demo",
		Title: "表单组件演示",
		Fields: []Field{
			textField("input", "单行文本", "blur"),
			textField("textarea", "多行文本", "blur"),
			textField("password", "密码", "blur"),
			choiceField("select", "下拉选择", "a", choice, "change"),
			choiceField("radio", "单选", "a", choice, "change"),
			{
				Type:    "checkbox",
				ID:      "checkbox",
				Title:   "多选",
				Default: []string{"a"},
				Options: choice,
				Rules:   []Rule{{Required: true, Message: "请填写多选", Trigger: "change"}},
			},
			{
				Type:    "switch",
				ID:      "switch",
				Title:   "开关",
				Default: false,
				Rules:   []Rule{{Required: true, Message: "请填写开关", Trigger: "change"}},
			},
			textField("date", "日期", "change"),
			textField("time", "时间", "change"),
			textField("datetime", "日期时间", "change"),
			labelOnly("txt", "文字说明"),
			labelOnly("line", "分割线"),
			labelOnly("year", "年份"),
			labelOnly("yearMonth", "年月"),
			labelOnly("file", "文件"),
			labelOnly("image", "图片"),
			labelOnly("color", "颜色"),
			labelOnly("icon", "图标"),
			labelOnly("area", "地区"),
			labelOnly("dropdownSearch", "远程搜索"),
			labelOnly("editor", "富文本"),
			labelOnly("markdown", "Markdown"),
			labelOnly("treeRadio", "树形单选"),
			labelOnly("treeCheckbox", "树形多选"),
			labelOnly("treeEdit", "树形编辑"),
			labelOnly("table", "表格"),
			labelOnly("custom", "自定义"),
			{Type: "hidden", ID: "hidden", Title: "隐藏域", Default: "secret"},
		},
	}
}

func textField(typ, title, trigger string) Field {
	return Field{
		Type:        typ,
		ID:          typ,
		Title:       title,
		Default:     "",
		Placeholder: "请输入" + title,
		Rules:       []Rule{{Required: true, Message: "请填写" + title, Trigger: trigger}},
	}
}

func choiceField(typ, title, def string, options []Option, trigger string) Field {
	return Field{
		Type:    typ,
		ID:      typ,
		Title:   title,
		Default: def,
		Options: options,
		Rules:   []Rule{{Required: true, Message: "请填写" + title, Trigger: trigger}},
	}
}

func labelOnly(typ, title string) Field {
	return Field{Type: typ, ID: typ, Title: title, Default: ""}
}
