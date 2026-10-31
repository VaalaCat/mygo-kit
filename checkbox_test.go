package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestCheckbox 验证选择框的公开交互。
func TestCheckbox(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"checkbox", func(t *testing.T) {
			on := false
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneDark()); Checkbox(c, &on, "Enable") }, 300, 120)
			tt.Click("Enable")
			if !on {
				t.Fatal("not checked")
			}
			tt.Key(0, ui.KeySpace)
			if on {
				t.Fatal("Space did not toggle")
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestCheckboxDisabled 验证已有焦点或展开状态的选择框禁用后不再修改值。
func TestCheckboxDisabled(t *testing.T) {
	// 1. 先建立交互状态，再用下一帧禁用，最后发送输入。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"checkbox", func(t *testing.T) {
			on := false
			disabled := false
			changes := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Checkbox(c, &on, "Check").Disabled(disabled).Changed() {
					changes++
				}
			}, 400, 200)
			tt.Click("Check")
			if !on {
				t.Fatal("setup")
			}
			disabled = true
			tt.Frame()
			before := changes
			tt.Key(0, ui.KeySpace)
			tt.Click("Check")
			if !on || changes != before {
				t.Fatalf("disabled change: value=%v changes=%d before=%d", on, changes, before)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
