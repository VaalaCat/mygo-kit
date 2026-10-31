package kit

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/egoist/mygo/ui/markdown/parser"
	"strings"
)

// MarkdownOpt 指定流式显示、共享选区和应用提供的图片及链接行为。
type MarkdownOpt struct {
	Streaming                bool
	ShareSelection           bool
	OnLink                   func(string)
	Image                    func(string) ImageResource
	OnImage                  func(string)
	CopyLabel, CodeCopyLabel string
	footnoteTarget           *int
}

type markdownKey struct{}
type markdownFootnoteKey struct{}

// Markdown 显示可选择的流式 Markdown，使用直角代码块、任务框和图片。
func Markdown(c *ui.Context, source string, opt MarkdownOpt) ui.Element {
	// 1. 选区归本组件或调用方的可选择容器，不建立文档业务副本。
	if opt.CopyLabel == "" {
		opt.CopyLabel = text(c, "copyMarkdown")
	}
	if opt.CodeCopyLabel == "" {
		opt.CodeCopyLabel = text(c, "copyCode")
	}
	t := GetTheme(c)
	root := ui.Column(c).Radius(0)
	if !opt.ShareSelection {
		root.Selectable()
		root.ContextMenu(func(menu *ui.Menu) {
			menu.EditItems()
			if menu.Item(opt.CopyLabel).Chosen() {
				c.WriteClipboard(source)
			}
		})
	}
	// 2. 缓存框架选项，并按当前主题更新纯样式。
	opt.footnoteTarget = ui.Local(root, markdownFootnoteKey{}, func() int { return -1 })
	options := ui.Local(root, markdownKey{}, func() markdown.Options { return markdown.Options{} })
	options.ShareSelection = true
	options.OnLink = opt.OnLink
	options.Reveal = &markdown.Reveal{Unit: markdown.AtOnce}
	th := markdown.DefaultTheme(c.Theme())
	th.Font = t.Typography.Font
	th.FontSize = t.Typography.BodySize
	th.CodeFont = t.Typography.MonoFont
	th.CodeSize = t.Typography.CodeSize / t.Typography.BodySize
	// 正文用 Prose 14/24、代码用 Code 12/20；块间距 8、紧凑列表项 4、代码块内边距 12，均换算为字号倍数。
	th.LineHeight = t.Space(6) / t.Typography.BodySize
	th.CodeLineHeight = t.Space(5) / t.Typography.CodeSize
	th.BlockSpacing = t.Space(2) / t.Typography.BodySize
	th.ItemSpacing = t.Space(1) / t.Typography.BodySize
	th.CodeBlockPadding = t.Space(3) / t.Typography.BodySize
	th.Text = t.Colors.Text
	th.Muted = t.Colors.Muted
	th.Link = t.Colors.Accent
	th.Rule = t.Colors.BorderSubtle
	th.TableBorder = t.Colors.BorderSubtle
	th.TableHeader = t.Colors.Hover
	th.CodeBackground = t.Colors.Hover
	th.CodeBlockBackground = t.Colors.Editor
	th.CodeBlockBorder = t.Colors.BorderSubtle
	th.CodeColors = codeTheme(c)
	options.Theme = th
	options.CodeBlock = func(c *ui.Context, block *markdown.CodeBlockView) {
		contentFrame(c, ui.Column(c)).Children(func() {
			frameHeader(c, block.Language, opt.CodeCopyLabel, block.Open, func() { c.WriteClipboard(block.Code()) })
			block.Body(c)
		})
	}
	// 3. 使用原文档节点绘制任务列表、图片和行内代码，不重新解析子片段。
	options.Block = func(c *ui.Context, d *parser.Document, index int) bool {
		block := d.Block(index)
		if opt.footnoteTarget != nil && block.Parent == *opt.footnoteTarget {
			ui.Box(c).Height(0).ScrollIntoView()
			*opt.footnoteTarget = -1
		}
		if block.Kind == parser.Table {
			markdownBlock(c, d, index, opt)
			return true
		}
		if block.Kind == parser.List {
			for item := range d.Children(index) {
				if d.Block(item).Task != parser.NoTask {
					markdownBlock(c, d, index, opt)
					return true
				}
			}
		}
		if block.Kind == parser.Paragraph || block.Kind == parser.Heading || block.Kind == parser.TableCell {
			for _, inline := range d.Inlines(index) {
				if inline.Kind == parser.Image || inline.Kind == parser.Code {
					markdownBlock(c, d, index, opt)
					return true
				}
			}
		}
		return false
	}
	root.Children(func() {
		if opt.Streaming {
			markdown.Stream(c, source, options)
		} else {
			markdown.View(c, source, options)
		}
	})
	return root
}

type markdownPiece struct {
	spans []ui.Span
	link  string
	image *parser.Inline
	alt   string
	note  int
}

// markdownParagraph 从同一 Document 的 Inline 树读取样式、链接和图片父级。
func markdownParagraph(c *ui.Context, d *parser.Document, index int, opt MarkdownOpt) {
	// 1. 扁平化显示片段时保留父级样式与链接目标。
	t := GetTheme(c)
	inlines := d.Inlines(index)
	pieces := []markdownPiece{}
	var visit func(int, int, ui.Span, string)
	visit = func(start, end int, style ui.Span, link string) {
		for i := start; i < end; {
			in := inlines[i]
			next := in.Next
			if next <= i || next > end {
				next = i + 1
			}
			switch in.Kind {
			case parser.Strong:
				child := style
				child.Weight = 600
				visit(i+1, next, child, link)
			case parser.Emphasis:
				child := style
				child.Italic = true
				visit(i+1, next, child, link)
			case parser.Strikethrough:
				child := style
				child.Strikethrough = true
				visit(i+1, next, child, link)
			case parser.Link:
				visit(i+1, next, style, in.Dest)
			case parser.Image:
				pieces = append(pieces, markdownPiece{link: link, image: &in, alt: parser.PlainText(inlines[i+1 : next])})
			default:
				span := style
				span.Text = in.Text
				note := 0
				switch in.Kind {
				case parser.SoftBreak:
					span.Text = " "
				case parser.LineBreak:
					span.Text = "\n"
				case parser.Code:
					span.Font = t.Typography.MonoFont
					span.Size = t.Typography.CodeSize
					span.Background = t.Colors.Hover
					span.Radius = t.Metrics.Radius
				case parser.FootnoteRef:
					if target, ok := d.Footnote(in.Text); ok {
						note = target
					}
					number := 0
					for j := 1; j <= note; j++ {
						if d.Block(j).Kind == parser.FootnoteDefinition {
							number++
						}
					}
					span.Text = fmt.Sprintf("[%d]", number)
					span.Color = t.Colors.Accent
					span.Size = t.Typography.CaptionSize
					span.BaselineShift = 3
				}
				if note == 0 && len(pieces) > 0 && pieces[len(pieces)-1].note == 0 && pieces[len(pieces)-1].image == nil && pieces[len(pieces)-1].link == link {
					pieces[len(pieces)-1].spans = append(pieces[len(pieces)-1].spans, span)
				} else {
					pieces = append(pieces, markdownPiece{spans: []ui.Span{span}, link: link, note: note})
				}
			}
			i = next
		}
	}
	style := ui.Span{Color: t.Colors.Text}
	if d.Block(index).Kind == parser.Heading {
		style.Size = t.Typography.HeadingSize
		style.Weight = 600
	}
	visit(0, len(inlines), style, "")
	// 2. 同一段的文字保留行内布局；图片另占行且保留外层链接操作。
	for i := 0; i < len(pieces); {
		piece := pieces[i]
		if piece.image != nil {
			resource := ImageResource{Alt: piece.alt}
			if opt.Image != nil {
				resource = opt.Image(piece.image.Dest)
				if resource.Alt == "" {
					resource.Alt = piece.alt
				}
			}
			if piece.link != "" || opt.OnImage != nil {
				button := ui.ButtonBase(c).Radius(t.Metrics.Radius).Label(piece.alt)
				if piece.link != "" {
					button.Role(ui.RoleLink)
				}
				button.Children(func() { Image(c, resource) })
				if button.Clicked() {
					if piece.link != "" {
						openMarkdownLink(c, piece.link, opt)
					} else {
						opt.OnImage(piece.image.Dest)
					}
				}
			} else {
				Image(c, resource)
			}
			i++
			continue
		}
		end := i
		for end < len(pieces) && pieces[end].image == nil {
			end++
		}
		ui.RichText(c).Children(func() {
			for _, part := range pieces[i:end] {
				if part.note > 0 {
					link := ui.Link(c, "", "")
					link.Children(func() { ui.RichText(c, part.spans...) })
					if link.Clicked() && opt.footnoteTarget != nil {
						*opt.footnoteTarget = part.note
					}
				} else if part.link == "" {
					ui.RichText(c, part.spans...)
				} else {
					url := part.link
					if opt.OnLink != nil || !ui.SafeToOpen(url) {
						url = ""
					}
					link := ui.Link(c, "", url).TextColor(t.Colors.Accent)
					link.Children(func() { ui.RichText(c, part.spans...) })
					if link.Clicked() && opt.OnLink != nil {
						opt.OnLink(part.link)
					}
				}
			}
		})
		i = end
	}
}

// openMarkdownLink 将链接交给调用方或系统安全外链入口。
func openMarkdownLink(c *ui.Context, url string, opt MarkdownOpt) {
	// 1. 只处理明确允许的外链，不执行任意协议。
	if opt.OnLink != nil {
		opt.OnLink(url)
	} else if ui.SafeToOpen(url) {
		c.Services().OpenURL(url)
	}
}

// markdownBlock 在原文档中递归构建任务项的后代，保留引用和脚注解析结果。
func markdownBlock(c *ui.Context, d *parser.Document, index int, opt MarkdownOpt) {
	// 1. 容器保留原节点关系，叶子直接使用已解析内容。
	t := GetTheme(c)
	block := d.Block(index)
	switch block.Kind {
	case parser.Paragraph, parser.Heading, parser.TableCell:
		markdownParagraph(c, d, index, opt)
	case parser.List:
		ui.Column(c).Role(ui.RoleList).Gap(t.Space(1)).Children(func() {
			number := block.Number
			for item := range d.Children(index) {
				task := d.Block(item).Task
				ui.Row(c).Role(ui.RoleListItem).Gap(t.Space(2)).Children(func() {
					if task != parser.NoTask {
						box := ui.Box(c).Size(checkBoxSize, checkBoxSize).Margin(t.Space(1), 0, 0, 0).Radius(t.Metrics.Radius).Role(ui.RoleCheckBox).Checked(task == parser.Checked).Label(text(c, "task")).Border(1, t.Colors.Muted)
						if task == parser.Checked {
							box.Background(t.Colors.Accent).Draw(func(p *ui.Painter, r ui.Rect) { checkMark(p, r, t.Colors.AccentText) })
						}
					} else {
						marker := "•"
						if block.Ordered {
							marker = fmt.Sprintf("%d.", number)
						}
						ui.Text(c, marker).TextColor(t.Colors.Muted)
					}
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						for child := range d.Children(item) {
							markdownBlock(c, d, child, opt)
						}
					})
				})
				number++
			}
		})
	case parser.CodeBlock:
		// 2. 代码文字不含文档引用，交给框架保留高亮和复制。
		code := d.Lines(index)
		longest, run := 0, 0
		for _, r := range code {
			if r == '`' {
				run++
				longest = max(longest, run)
			} else {
				run = 0
			}
		}
		fence := strings.Repeat("`", max(3, longest+1))
		text := fence + block.Language() + "\n" + code
		if !opt.Streaming || !block.Open {
			text += fence
		}
		nested := opt
		nested.ShareSelection = true
		nested.Streaming = opt.Streaming && block.Open
		Markdown(c, text, nested)
	case parser.Quote:
		ui.Column(c).BorderStartWidth(2).BorderColor(t.Colors.Border).PaddingStart(t.Space(3)).Children(func() {
			for child := range d.Children(index) {
				markdownBlock(c, d, child, opt)
			}
		})
	case parser.ThematicBreak:
		Separator(c)
	case parser.HTMLBlock, parser.MathBlock:
		t.TextRole(ui.Text(c, d.Lines(index)), TextCode)
	case parser.Table:
		widths := make([]float32, len(block.Align))
		for row := range d.Children(index) {
			col := 0
			for cell := range d.Children(row) {
				if col < len(widths) {
					weight := 400
					if d.Block(row).Header {
						weight = 600
					}
					w, _ := c.MeasureText(0, ui.Span{Text: parser.PlainText(d.Inlines(cell)), Size: t.Typography.BodySize, Weight: weight})
					widths[col] = max(widths[col], min(360, w)+24)
				}
				col++
			}
		}
		ui.ScrollHorizontal(c).Children(func() {
			ui.Column(c).Radius(t.Metrics.Radius).Border(1, t.Colors.BorderSubtle).Clip().Shrink(0).Role(ui.RoleTable).Children(func() {
				for row := range d.Children(index) {
					line := ui.Row(c).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.BorderSubtle).Role(ui.RoleRow)
					if d.Block(row).Header {
						line.Background(t.Colors.Hover).FontWeight(600)
					}
					line.Children(func() {
						col := 0
						for cell := range d.Children(row) {
							role := ui.RoleCell
							if d.Block(row).Header {
								role = ui.RoleColumnHeader
							}
							box := ui.Column(c).Width(widths[col]).Padding(t.Space(1), t.Space(2)).Role(role)
							if col > 0 {
								box.BorderStartWidth(1).BorderColor(t.Colors.BorderSubtle)
							}
							box.Children(func() { markdownParagraph(c, d, cell, opt) })
							col++
						}
					})
				}
			})
		})
	default:
		for child := range d.Children(index) {
			markdownBlock(c, d, child, opt)
		}
	}
}
