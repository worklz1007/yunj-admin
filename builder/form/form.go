package form

type Form struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Default     any      `json:"default"`
	Desc        string   `json:"desc,omitempty"`
	Readonly    bool     `json:"readonly"`
	Placeholder string   `json:"placeholder,omitempty"`
	Options     []Option `json:"options,omitempty"`
	Rules       []Rule   `json:"rules,omitempty"`
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Rule struct {
	Required bool   `json:"required"`
	Message  string `json:"message"`
	Trigger  string `json:"trigger"`
}
