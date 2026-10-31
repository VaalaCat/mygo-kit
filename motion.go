package kit

import (
	"github.com/egoist/mygo/ui"
	"github.com/go-text/typesetting/segmenter"
	"strings"
	"time"
)

// Fade 更新元素可见透明度；系统或主题请求减少动画时立即到达目标。
func Fade(c *ui.Context, element ui.Element, visible bool) ui.Element {
	// 1. 隐藏中的内容立即退出交互，透明度只影响绘制。
	target := float32(0)
	if visible {
		target = 1
	}
	t := GetTheme(c)
	if !visible {
		element.Disabled(true)
	}
	if t.Motion.Reduced || c.Preferences().ReduceMotion {
		return element.Opacity(target)
	}
	return element.Opacity(element.Animate("kit-fade", target, t.Motion.Fade))
}

type morphKey struct{}
type morphGlyph struct {
	text string
	id   uint64
	born time.Time
}
type morphTextState struct {
	text   string
	glyphs []morphGlyph
	nextID uint64
}

// MorphText 显示单行状态文字，只对变化字素短上推并渐显；共同前后缀及中段匹配保持身份，纯追加不动画。
func MorphText(c *ui.Context, value string) ui.Element {
	// 1. 按字素更新有限匹配结果，不拆分组合字符。
	theme := GetTheme(c)
	root := ui.Row(c).Gap(0).AlignItems(ui.Center)
	state := ui.Local(root, morphKey{}, func() morphTextState { return morphTextState{} })
	if state.glyphs == nil || state.text != value {
		state.update(value, c.Now())
	}
	reduced := theme.Motion.Reduced || c.Preferences().ReduceMotion

	// 2. 相同变化时刻的连续字素组成文本 run；节点透明度支持真正的 0。
	root.Children(func() {
		for start := 0; start < len(state.glyphs); {
			glyph := state.glyphs[start]
			end := start + 1
			for end < len(state.glyphs) && state.glyphs[end].born.Equal(glyph.born) {
				end++
			}
			var content strings.Builder
			for _, part := range state.glyphs[start:end] {
				content.WriteString(part.text)
			}
			element := theme.TextRole(ui.Text(c.Key(glyph.id), content.String()).NoWrap().Shrink(0), TextUI)
			if !reduced && !glyph.born.IsZero() {
				elapsed := c.Now().Sub(glyph.born)
				duration := theme.Motion.Fade
				if duration > 0 && elapsed < duration {
					alpha := max(0, min(1, float32(elapsed)/float32(duration)))
					element.Opacity(alpha).Translate(0, 3*(1-alpha))
					c.AnimationFrame()
				}
			}
			start = end
		}
	})

	return root.Label(value)
}

// update 对短变化段做有界匹配，超大变化直接显示完整文本。
func (s *morphTextState) update(value string, now time.Time) {
	// 1. 先保留共同前缀和后缀，纯追加的新字符不入场。
	next := textGraphemes(value)
	old := s.glyphs
	prefix := 0
	for prefix < len(old) && prefix < len(next) && old[prefix].text == next[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(old)-prefix && suffix < len(next)-prefix && old[len(old)-suffix-1].text == next[len(next)-suffix-1] {
		suffix++
	}
	result := make([]morphGlyph, len(next))
	copy(result[:prefix], old[:prefix])
	if suffix > 0 {
		copy(result[len(next)-suffix:], old[len(old)-suffix:])
	}
	oldEnd, newEnd := len(old)-suffix, len(next)-suffix
	animate := len(old) > 0 && prefix != len(old)
	create := func(index int) {
		s.nextID++
		born := time.Time{}
		if animate {
			born = now
		}
		result[index] = morphGlyph{text: next[index], id: s.nextID, born: born}
	}
	// 2. 变化中段最多比较 16384 对字素，保留其中未变化字符的身份。
	n, m := oldEnd-prefix, newEnd-prefix
	if n*m > 16384 || n > 256 || m > 256 {
		animate = false
		for i := prefix; i < newEnd; i++ {
			create(i)
		}
	} else {
		lengths := make([]int, (n+1)*(m+1))
		at := func(i, j int) int { return i*(m+1) + j }
		for i := n - 1; i >= 0; i-- {
			for j := m - 1; j >= 0; j-- {
				if old[prefix+i].text == next[prefix+j] {
					lengths[at(i, j)] = lengths[at(i+1, j+1)] + 1
				} else {
					lengths[at(i, j)] = max(lengths[at(i+1, j)], lengths[at(i, j+1)])
				}
			}
		}
		i, j := 0, 0
		for j < m {
			if i < n && old[prefix+i].text == next[prefix+j] {
				result[prefix+j] = old[prefix+i]
				i++
				j++
			} else if i < n && lengths[at(i+1, j)] >= lengths[at(i, j+1)] {
				i++
			} else {
				create(prefix + j)
				j++
			}
		}
	}
	s.text = value
	s.glyphs = result
}

// textGraphemes 按 Unicode 字素分割文本。
func textGraphemes(value string) []string {
	// 1. 完整 emoji、组合重音和文字簇作为动画单位。
	var segments segmenter.Segmenter
	segments.Init([]rune(value))
	iterator := segments.GraphemeIterator()
	var result []string
	for iterator.Next() {
		result = append(result, string(iterator.Grapheme().Text))
	}
	return result
}
