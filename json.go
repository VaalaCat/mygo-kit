package kit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/egoist/mygo/ui"
	"io"
	"slices"
	"strconv"
	"strings"
)

// JSONState 保存已解析的 JSON 与展开、选中、滚动状态。
type JSONState struct {
	list     ui.ListState
	selected int
	// open 是已展开节点的路径集合，根节点始终在其中直到被用户收起。
	open  map[string]bool
	nodes map[string]jsonNode
	// rows 是当前可见节点按显示顺序的展平结果，展开状态或内容变化后重建。
	rows      []jsonRow
	rowsStale bool
}

type jsonNode struct {
	label    string
	value    any
	children []string
}

// jsonRow 是展平后的一行：节点路径、缩进层级和父行下标（根为 -1）。
type jsonRow struct {
	path   string
	depth  int
	parent int
}

// NewJSONState 解析任意 JSON 顶级值，包括 null、数组和标量。
func NewJSONState(data []byte) (*JSONState, error) {
	// 1. 创建树并解析输入。
	s := &JSONState{selected: -1, open: map[string]bool{}}
	if err := s.Set(data); err != nil {
		return nil, err
	}
	return s, nil
}

// Set 解析新内容，成功后替换树，保留仍存在的展开路径；根节点只在第一次设置时展开。
func (s *JSONState) Set(data []byte) error {
	// 1. 完整校验一个 JSON 值，不将浮点数替代原数字。
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	// 2. 建立稳定的 JSON Pointer 路径。
	nodes := map[string]jsonNode{}
	var build func(string, string, any)
	build = func(path, label string, value any) {
		n := jsonNode{label: label, value: value}
		switch v := value.(type) {
		case map[string]any:
			names := make([]string, 0, len(v))
			for name := range v {
				names = append(names, name)
			}
			slices.Sort(names)
			for _, name := range names {
				child := path + "/" + strings.ReplaceAll(strings.ReplaceAll(name, "~", "~0"), "/", "~1")
				n.children = append(n.children, child)
				build(child, name, v[name])
			}
		case []any:
			for i, item := range v {
				child := path + "/" + strconv.Itoa(i)
				n.children = append(n.children, child)
				build(child, strconv.Itoa(i), item)
			}
		}
		nodes[path] = n
	}
	build("$", "$", value)
	// 3. 只在第一次设置时展开根节点，之后用户收起的根节点不会被刷新内容重新展开。
	if s.nodes == nil {
		s.open["$"] = true
	}
	s.nodes = nodes
	s.rowsStale = true
	s.selected = -1
	return nil
}

// flatten 按展开状态把可见节点展平成行，行内容没有变化时沿用上一次结果。
func (s *JSONState) flatten() {
	if !s.rowsStale {
		return
	}
	s.rows = s.rows[:0]
	var walk func(path string, depth, parent int)
	walk = func(path string, depth, parent int) {
		s.rows = append(s.rows, jsonRow{path: path, depth: depth, parent: parent})
		if !s.open[path] {
			return
		}
		self := len(s.rows) - 1
		for _, child := range s.nodes[path].children {
			walk(child, depth+1, self)
		}
	}
	walk("$", 0, -1)
	s.rowsStale = false
}

// setOpen 展开或收起一个节点，并让下一帧重新展平。
func (s *JSONState) setOpen(path string, open bool) {
	if open {
		s.open[path] = true
	} else {
		delete(s.open, path)
	}
	s.rowsStale = true
}

// choose 选中第 i 行并滚动到可见处。
func (s *JSONState) choose(i int) {
	s.selected = i
	s.list.ScrollIntoView(i)
}

// handleKeys 处理选中节点的展开、收起和跳转：Right 展开或进入第一个子节点，Left 收起或回到父节点，Enter 与 Space 切换。
// rowClicked 表示本帧已有行被点击切换过，列表把双击也报告为 Submitted，此时不再重复切换。
func (s *JSONState) handleKeys(list ui.Element, rowClicked bool) {
	// 1. 每个按键都向框架登记，不能因前一个命中而短路；Enter 被列表先行处理，只能读它的 Submitted。
	right, left := list.Shortcut(0, ui.KeyRight), list.Shortcut(0, ui.KeyLeft)
	enter, space := list.Submitted() && !rowClicked, list.Shortcut(0, ui.KeySpace)
	if s.selected < 0 || s.selected >= len(s.rows) {
		return
	}
	// 2. 只有含成员的节点可以展开；展开后第一个子节点紧随其后。
	row := s.rows[s.selected]
	branch := len(s.nodes[row.path].children) > 0
	open := s.open[row.path]
	switch {
	case right && branch && !open:
		s.setOpen(row.path, true)
	case right && branch && open:
		s.choose(s.selected + 1)
	case left && branch && open:
		s.setOpen(row.path, false)
	case left && row.parent >= 0:
		s.choose(row.parent)
	case (enter || space) && branch:
		s.setOpen(row.path, !open)
	}
}

// JSONOpt 设置 JSON 工具栏标题及复制按钮的显示和可访问名称。
type JSONOpt struct{ Title, CopyLabel string }

// jsonListMinHeight 是 JSON 树列表的最小高度。
const jsonListMinHeight float32 = 48

// JSONTree 显示可键盘展开、选择和复制的虚拟 JSON 树。
func JSONTree(c *ui.Context, s *JSONState, opt JSONOpt) ui.Element {
	// 1. 标题与复制名称为空时采用默认值；选择与复制直接使用同一棵解析树。
	if opt.Title == "" {
		opt.Title = "JSON"
	}
	if opt.CopyLabel == "" {
		opt.CopyLabel = text(c, "copyJSON")
	}
	t := GetTheme(c)
	root := contentFrame(c, ui.Column(c).Gap(0))
	root.Children(func() {
		frameHeader(c, opt.Title, opt.CopyLabel, false, func() { c.WriteClipboard(jsonText(s.nodes["$"].value)) })
		// 2. 虚拟列表只构建可见行；行路径作为身份，展开或收起时选中项跟着节点走。
		s.flatten()
		s.list.Selected = &s.selected
		list := ui.List(c, &s.list, len(s.rows)).Grow(1).MinHeight(jsonListMinHeight).Radius(0).Gap(0).Padding(0).ItemKey(func(i int) any { return s.rows[i].path })
		rowClicked := false
		list.Rows(func(row ui.ListRow) { rowClicked = s.buildRow(c, t, row) || rowClicked })
		s.handleKeys(list, rowClicked)
		if list.FocusWithin() && list.Shortcut(ui.Cmd, ui.KeyC) && s.selected >= 0 && s.selected < len(s.rows) {
			c.WriteClipboard(jsonText(s.nodes[s.rows[s.selected].path].value))
		}
	})
	return root
}

// buildRow 绘制一行：20 高，缩进每级 16，其后是 16 宽的展开箭头位、键名和值；点击整行选中并切换展开，返回本行是否被点击。
func (s *JSONState) buildRow(c *ui.Context, t Theme, listRow ui.ListRow) (clicked bool) {
	// 1. 底色只在选中时给出，悬停底色由列表行自己绘制。
	i := listRow.Index
	row, n := s.rows[i], s.nodes[s.rows[i].path]
	branch := len(n.children) > 0
	open := s.open[row.path]
	line := ui.Row(c).FillWidth().Height(t.Metrics.CompactHeight).Shrink(0).Gap(0).AlignItems(ui.Center).Padding(0, t.Space(3), 0, t.Space(3)+float32(row.depth)*t.Space(4))
	line.Background(controlBackground(&t, controlState{Selected: listRow.Selected()}))
	if line.Clicked() {
		clicked = true
		s.selected = i
		if branch {
			s.setOpen(row.path, !open)
		}
	}
	line.Children(func() {
		// 2. 箭头位对叶子留空，保持键名对齐。
		slot := ui.Row(c).Size(IconSizeLarge, IconSizeLarge).Shrink(0).Center()
		slot.Children(func() {
			if !branch {
				return
			}
			chevron := IconChevronRight
			if open {
				chevron = IconChevronDown
			}
			ui.Icon(c, chevron).Size(IconSizeSmall, IconSizeSmall).Shrink(0).TextColor(t.Colors.Muted)
		})
		// 3. 键名与值按基线对齐。
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Baseline).Grow(1).Shrink(1).MinWidth(0).Children(func() {
			t.TextRole(ui.Text(c, n.label).SingleLine().Shrink(0), TextCode).TextColor(t.Syntax.Property)
			// 4. 对象和数组只显示成员数，用次要色与值区分。
			switch v := n.value.(type) {
			case map[string]any:
				t.TextRole(ui.Text(c, fmt.Sprintf("{%d}", len(v))).SingleLine(), TextCode).TextColor(t.Colors.Muted)
			case []any:
				t.TextRole(ui.Text(c, fmt.Sprintf("[%d]", len(v))).SingleLine(), TextCode).TextColor(t.Colors.Muted)
			default:
				t.TextRole(ui.Text(c, jsonText(v)).SingleLine().Shrink(1).MinWidth(0), TextCode).TextColor(t.Colors.Text)
			}
		})
	})
	return clicked
}

// jsonText 序列化单个值，供显示和复制。
func jsonText(value any) string {
	// 1. 数字保留其 JSON 表达，&、<、> 不转义，显示和复制与真实值一致。
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
	return strings.TrimSuffix(out.String(), "\n")
}
