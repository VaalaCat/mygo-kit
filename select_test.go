package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestSelect 验证选择器的公开交互。
func TestSelect(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"select-local-value", func(t *testing.T) {
			value := "One"
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				local := value
				if Select(c, &local, SelectOpt{Options: []SelectOption{{Value: "One", Label: "One"}, {Value: "Two", Label: "Two"}}}).Label("Choice").Changed() {
					value = local
				}
			}, 400, 240)
			tt.Click("Choice")
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if value != "Two" {
				t.Fatalf("choice=%q", value)
			}
		}},
		{"select-label-placeholder-disabled", func(t *testing.T) {
			value := ""
			options := SelectOpt{Label: "Model", Placeholder: "Choose model", Options: []SelectOption{{Value: "a", Label: "Alpha"}, {Value: "blocked", Label: "Unavailable", Disabled: true}, {Value: "b", Label: "Beta"}}}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Select(c, &value, options) }, 400, 280)
			if !tt.HasText("Choose model") {
				t.Fatal(tt.Texts())
			}
			tt.Click("Model")
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if value != "a" {
				t.Fatal(value)
			}
			if !tt.HasText("Alpha") {
				t.Fatal(tt.Texts())
			}
			tt.Click("Model")
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if value != "b" {
				t.Fatalf("disabled value entered navigation: %q", value)
			}
		}},
		{"header-is-skipped-and-description-and-trigger-label", func(t *testing.T) {
			value := ""
			options := SelectOpt{Label: "Model", Options: []SelectOption{
				{Label: "Anthropic", Header: true},
				{Value: "a", Label: "Alpha", Description: "Provider A"},
				{Label: "OpenAI", Header: true},
				{Value: "b", Label: "Beta", Description: "Provider B"},
			}, TriggerLabel: func(selected SelectOption) string { return "Pick: " + selected.Label }}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Select(c, &value, options) }, 400, 280)
			tt.Click("Model")
			for _, want := range []string{"Anthropic", "OpenAI", "Provider A", "Provider B"} {
				if !tt.HasText(want) {
					t.Fatalf("missing %q in %q", want, tt.Texts())
				}
			}
			// 方向键第一次落在 Alpha（跳过标题），再按一次落在 Beta（跳过中间的标题）。
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if value != "b" {
				t.Fatalf("value=%q", value)
			}
			if !tt.HasText("Pick: Beta") || tt.HasText("Provider B") {
				t.Fatalf("trigger texts=%q", tt.Texts())
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestSelectDisabled 验证已有焦点或展开状态的选择器禁用后不再修改值。
func TestSelectDisabled(t *testing.T) {
	// 1. 先建立交互状态，再用下一帧禁用，最后发送输入。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"select-open", func(t *testing.T) {
			value := "a"
			disabled := false
			changes := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Select(c, &value, SelectOpt{Label: "Choice", Disabled: disabled, Options: []SelectOption{{Value: "a", Label: "Alpha"}, {Value: "b", Label: "Beta"}}}).Changed() {
					changes++
				}
			}, 400, 280)
			tt.Click("Choice")
			disabled = true
			tt.Frame()
			before := changes
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if value != "a" || changes != before {
				t.Fatalf("disabled keyboard change: value=%v changes=%d before=%d", value, changes, before)
			}
			tt.Click("Beta")
			if value != "a" || changes != before {
				t.Fatalf("disabled popup click changed value=%v", value)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestSelectSearch 验证搜索只出现在弹层内，筛选结果保留真实值，Enter、方向键和 Esc 在搜索框里的行为。
func TestSelectSearch(t *testing.T) {
	// 1. 每个用例打开弹层、输入过滤词后按给定的键，检查最终选中的值和弹层是否仍然打开。
	options := []SelectOption{{Value: "provider/a", Label: "Alpha"}, {Value: "provider/b", Label: "Beta"}, {Value: "provider/c", Label: "Bravo"}}
	tests := []struct {
		name     string
		typed    string
		keys     []ui.Key
		want     string
		wantOpen bool
		changed  bool
	}{
		{"enter-selects-first-match", "Bet", []ui.Key{ui.KeyEnter}, "provider/b", false, true},
		{"enter-without-filter-selects-first-option", "", []ui.Key{ui.KeyEnter}, "provider/a", false, true},
		{"down-then-enter-selects-second-match", "br", []ui.Key{ui.KeyDown, ui.KeyEnter}, "provider/c", false, true},
		{"down-stops-at-last-match", "b", []ui.Key{ui.KeyDown, ui.KeyDown, ui.KeyDown, ui.KeyEnter}, "provider/c", false, true},
		{"up-moves-back-to-first-match", "b", []ui.Key{ui.KeyDown, ui.KeyUp, ui.KeyEnter}, "provider/b", false, true},
		{"enter-without-match-keeps-popup", "zzz", []ui.Key{ui.KeyEnter}, "", true, false},
		{"escape-closes-without-choosing", "Bet", []ui.Key{ui.KeyEscape}, "", false, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selected, changed := "", false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Select(c, &selected, SelectOpt{Label: "Model", Searchable: true, SearchPlaceholder: "Find model", Options: options}).Changed() {
					changed = true
				}
			}, 400, 300)
			if tt.HasText("Find model") {
				t.Fatal("search leaked outside popup")
			}
			tt.Click("Model")
			if test.typed != "" {
				tt.Click("Find model")
				tt.Type(test.typed)
			}
			for _, key := range test.keys {
				tt.Key(0, key)
			}
			tt.Frame()
			if selected != test.want || changed != test.changed || tt.HasText("Find model") != test.wantOpen {
				t.Fatalf("selected=%q changed=%v open=%v, want %q %v %v", selected, changed, tt.HasText("Find model"), test.want, test.changed, test.wantOpen)
			}
		})
	}
}

// TestSelectVisibleOptions 验证搜索后分组标题只在组内有匹配项时保留，说明也参与匹配。
func TestSelectVisibleOptions(t *testing.T) {
	options := []SelectOption{
		{Label: "G1", Header: true}, {Value: "a", Label: "Alpha", Description: "red"}, {Value: "b", Label: "Beta"},
		{Label: "G2", Header: true}, {Value: "c", Label: "Gamma", Description: "blue"},
		{Label: "G3", Header: true},
	}
	tests := []struct {
		name       string
		query      string
		searchable bool
		want       []string
	}{
		{"no-search-keeps-everything", "", true, []string{"G1", "Alpha", "Beta", "G2", "Gamma", "G3"}},
		{"not-searchable-ignores-query", "zzz", false, []string{"G1", "Alpha", "Beta", "G2", "Gamma", "G3"}},
		{"label-match-drops-empty-groups", "beta", true, []string{"G1", "Beta"}},
		{"description-matches", "blue", true, []string{"G2", "Gamma"}},
		{"both-groups", "a", true, []string{"G1", "Alpha", "Beta", "G2", "Gamma"}},
		{"nothing-matches", "zzz", true, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got []string
			for _, option := range selectVisibleOptions(options, &test.query, test.searchable) {
				got = append(got, option.Label)
			}
			if len(got) != len(test.want) {
				t.Fatalf("got %v want %v", got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Fatalf("got %v want %v", got, test.want)
				}
			}
		})
	}
}
