package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type ProfanityExtension struct{}

func (e *ProfanityExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewProfanityParser(), 100),
			util.Prioritized(NewTooltipParser(), 101),
		),
	)
}

type ProfanityNode struct {
	ast.BaseInline
	Explicit string
	Safe     string
}

// Kind returns the kind of the node
func (n *ProfanityNode) Kind() ast.NodeKind {
	return ast.KindInline
}

// Type returns profanity type
func (n *ProfanityNode) Type() ast.NodeType {
	return ast.NodeInline
}

func (n *ProfanityNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

type TooltipNode struct {
	ast.BaseInline
	Text    string
	Tooltip string
}

// Kind returns the kind of the node
func (n *TooltipNode) Kind() ast.NodeKind {
	return ast.KindInline
}

// Type returns tooltip type
func (n *TooltipNode) Type() ast.NodeType {
	return ast.NodeInline
}

func (n *TooltipNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

func NewProfanityParser() parser.InlineParser {
	return &profanityParser{}
}

type profanityParser struct{}

func (p *profanityParser) Trigger() []byte { return []byte{'<'} }

func (p *profanityParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte("<prof>")) {
		return nil
	}

	// Parse profanity content
	block.Advance(6) // skip <prof>
	explicit := ""
	safe := ""

	// TODO: Implement full parsing logic for <prof> and <safe> tags

	node := &ProfanityNode{
		Explicit: explicit,
		Safe:     safe,
	}

	return node
}

func NewTooltipParser() parser.InlineParser {
	return &tooltipParser{}
}

type tooltipParser struct{}

func (p *tooltipParser) Trigger() []byte { return []byte{'<'} }

func (p *tooltipParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, []byte("<tt ")) {
		return nil
	}

	// TODO: Implement tooltip parsing logic

	node := &TooltipNode{
		Text:    "placeholder",
		Tooltip: "tooltip text",
	}

	return node
}

func ProcessMarkdown(content []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(&ProfanityExtension{}),
	)

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
