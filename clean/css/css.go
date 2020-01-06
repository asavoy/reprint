package css

import (
	"path"
	"regexp"
	"strings"

	"github.com/vanng822/css"
)

var (
	urlRegexp = regexp.MustCompile(`url\("(.+?)"\)|url\('(.+?)'\)|url\((.+?)\)`)
)

func RebaseURLs(ss *css.CSSStyleSheet, curPath string, targetPath string) {
	curPath = path.Clean(curPath)
	targetPath = path.Clean(targetPath)

	relPathToRoot := strings.Repeat("../", len(strings.Split(targetPath, "/")))
	relPathToCur := path.Clean(path.Join(relPathToRoot, curPath))

	for _, rule := range ss.CssRuleList {
		rebaseURLs(rule, relPathToCur)
	}
}

func rebaseURLs(rule *css.CSSRule, relPathToCur string) {
	for _, childRule := range rule.Rules {
		rebaseURLs(childRule, relPathToCur)
	}
	var newStyles []*css.CSSStyleDeclaration
	for _, style := range rule.Style.Styles {
		style.Value = urlRegexp.ReplaceAllStringFunc(style.Value, func(match string) string {
			matches := urlRegexp.FindStringSubmatch(match)
			url := matches[1] + matches[2] + matches[3]
			newURL := path.Clean(path.Join(relPathToCur, url))
			return `url("` + newURL + `")`
		})
		newStyles = append(newStyles, style)
	}
	rule.Style.Styles = newStyles
}

func RemoveMediaRules(ss *css.CSSStyleSheet) {
	var rules []*css.CSSRule
	for _, rule := range ss.CssRuleList {
		if rule.Type != css.MEDIA_RULE {
			rules = append(rules, rule)
		}
	}
	ss.CssRuleList = rules
}

func RemoveKeyframeRules(ss *css.CSSStyleSheet) {
	var rules []*css.CSSRule
	for _, rule := range ss.CssRuleList {
		if rule.Type != css.KEYFRAMES_RULE {
			rules = append(rules, rule)
		}
	}
	ss.CssRuleList = rules
}

func RemoveFontFaceRules(ss *css.CSSStyleSheet) {
	var rules []*css.CSSRule
	for _, rule := range ss.CssRuleList {
		if rule.Type != css.FONT_FACE_RULE {
			rules = append(rules, rule)
		}
	}
	ss.CssRuleList = rules
}

func RemoveColors(ss *css.CSSStyleSheet) {
	for _, rule := range ss.CssRuleList {
		noColors(rule)
	}
}

func noColors(rule *css.CSSRule) {
	for _, childRule := range rule.Rules {
		noColors(childRule)
	}
	var newStyles []*css.CSSStyleDeclaration
	for _, style := range rule.Style.Styles {
		if style.Property == "color" || style.Property == "background-color" {
			continue
		} else {
			newStyles = append(newStyles, style)
		}
	}
	rule.Style.Styles = newStyles
}

func RemoveTextAlignJustify(ss *css.CSSStyleSheet) {
	for _, rule := range ss.CssRuleList {
		noTextAlignJustify(rule)
	}
}

func noTextAlignJustify(rule *css.CSSRule) {
	for _, childRule := range rule.Rules {
		noTextAlignJustify(childRule)
	}
	var newStyles []*css.CSSStyleDeclaration
	for _, style := range rule.Style.Styles {
		if style.Property == "text-align" && style.Value == "justify" {
			continue
		} else {
			newStyles = append(newStyles, style)
		}
	}
	rule.Style.Styles = newStyles
}

func KeepSimpleStyles(ss *css.CSSStyleSheet) {
	for _, rule := range ss.CssRuleList {
		simpleStyles(rule)
	}
}

func simpleStyles(rule *css.CSSRule) {
	keepStyles := map[string]bool{
		"background-color": true,
		"color":            true,
		"content":          true,
		"display":          true,
		"font-style":       true,
		"font-weight":      true,
		"font-decoration":  true,
		"text-align":       true,
		"text-transform":   true,
		"white-space":      true,
	}
	for _, childRule := range rule.Rules {
		simpleStyles(childRule)
	}
	var newStyles []*css.CSSStyleDeclaration
	for _, style := range rule.Style.Styles {
		if _, ok := keepStyles[style.Property]; ok {
			newStyles = append(newStyles, style)
		}
	}
	rule.Style.Styles = newStyles
}

func AddHeadingStyles(ss *css.CSSStyleSheet) {
	ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: "h1, h2, h3, h4, h5, h6",
			Styles: []*css.CSSStyleDeclaration{
				{Property: "font-weight", Value: "bold"},
				{Property: "-webkit-hyphens", Value: "none", Important: 1},
				{Property: "hyphens", Value: "none", Important: 1},
				{Property: "page-break-inside", Value: "avoid"},
				{Property: "page-break-after", Value: "avoid"},
			},
		},
		Rules: nil,
	})

	// To distinguish smaller headings from body text
	ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: "h5, h6",
			Styles: []*css.CSSStyleDeclaration{
				{Property: "text-transform", Value: "uppercase"},
			},
		},
		Rules: nil,
	})
}

func AddFigureStyles(ss *css.CSSStyleSheet) {
	ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: ".figure, figure",
			Styles: []*css.CSSStyleDeclaration{
				{Property: "page-break-inside", Value: "avoid"},
			},
		},
		Rules: nil,
	})
}

func AddAsideStyles(ss *css.CSSStyleSheet) {
	ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: `aside, .aside, .box, .boxg, .note, .note1, sidebar, .sidebar1, [data-type="note"], [data-type="tip"], [data-type="warning"]`,
			Styles: []*css.CSSStyleDeclaration{
				{Property: "border", Value: "1px dotted #ddd"},
				{Property: "padding", Value: "0em 1em", Important: 1},
				{Property: "margin-top", Value: "1em", Important: 1},
				{Property: "margin-bottom", Value: "1em", Important: 1},
				{Property: "page-break-inside", Value: "avoid"},
			},
		},
		Rules: nil,
	})
}

func AddTableStyles(ss *css.CSSStyleSheet) {
	ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: `table`,
			Styles: []*css.CSSStyleDeclaration{
				{Property: "border-collapse", Value: "collapse"},
			},
		},
		Rules: nil,
	}, &css.CSSRule{
		Type: css.STYLE_RULE,
		Style: css.CSSStyleRule{
			SelectorText: `td, th`,
			Styles: []*css.CSSStyleDeclaration{
				{Property: "padding", Value: "0 0.5em"},
			},
		},
		Rules: nil,
	})
}

func Render(ss *css.CSSStyleSheet) string {
	rendered := ""
	for _, rule := range ss.GetCSSRuleList() {
		rendered += renderRule(rule, 0)
	}
	return rendered
}

func renderRule(rule *css.CSSRule, indentLevel int) string {
	indent := "    "
	baseIndent := strings.Repeat(indent, indentLevel)
	childIndent := baseIndent + indent

	children := ""

	for _, style := range rule.Style.Styles {
		children += childIndent + style.Text() + ";\n"
	}
	for _, childRule := range rule.Rules {
		children += renderRule(childRule, indentLevel+1)
	}
	if children == "" {
		return ""
	}
	rendered := baseIndent
	if rule.Type.Text() != "" {
		rendered += rule.Type.Text() + " "
	}
	if rule.Style.SelectorText != "" {
		rendered += rule.Style.SelectorText + " "
	}
	rendered += "{\n" + children + baseIndent + "}\n"
	return rendered
}
