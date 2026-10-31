package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestSheet 验证抽屉由按钮打开，Escape 或点遮罩关闭，并贴靠指定一侧。
func TestSheet(t *testing.T) {
	// 1. 每个用例使用一个打开按钮和一个抽屉。
	tests := []struct {
		name  string
		side  SheetSide
		close func(*ui.Tester)
		left  bool
	}{
		{"left-escape", SheetLeft, func(tt *ui.Tester) { tt.Key(0, ui.KeyEscape) }, true},
		{"right-escape", SheetRight, func(tt *ui.Tester) { tt.Key(0, ui.KeyEscape) }, false},
		{"right-backdrop", SheetRight, func(tt *ui.Tester) { tt.ClickAt(20, 300) }, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			open := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Button(c, "Open").Clicked() {
					open = true
				}
				Sheet(c, &open, test.side, 300, func() { ui.Text(c, "Sheet body") })
			}, 800, 600)
			if tt.HasText("Sheet body") {
				t.Fatal("open before click")
			}
			tt.Click("Open")
			body, ok := tt.Find("Sheet body")
			if !open || !ok {
				t.Fatal("sheet did not open")
			}
			if test.left && body.X > 300 || !test.left && body.X < 800-300 {
				t.Fatalf("sheet is on the wrong side: %v", body)
			}
			test.close(tt)
			if open || tt.HasText("Sheet body") {
				t.Fatal("sheet did not close")
			}
		})
	}
}

// TestSheetWithoutTitleBar 验证窗口没有 TitleBar 时抽屉和遮罩铺满整个窗口，顶部点击落在遮罩上。
func TestSheetWithoutTitleBar(t *testing.T) {
	// 1. 抽屉打开后点窗口最上方的一点，遮罩接住点击并关闭抽屉。
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		Sheet(c, &open, SheetRight, 300, func() { ui.Text(c, "Sheet body") })
	}, 800, 600)
	tt.Frame()
	if body, _ := tt.Find("Sheet body"); body.Y >= 36 {
		t.Fatalf("sheet reserves a title bar band that does not exist: %v", body)
	}
	tt.ClickAt(100, 4)
	if open {
		t.Fatal("click on the top strip did not hit the scrim")
	}
}

// TestSheetBelowTitleBar 验证抽屉和遮罩从顶栏下缘开始：顶栏按钮仍能点击，点遮罩才关闭。
func TestSheetBelowTitleBar(t *testing.T) {
	// 1. 顶栏占窗口最上方一行，抽屉打开后分别点顶栏按钮和遮罩。
	theme := OneLight()
	tests := []struct {
		name        string
		side        SheetSide
		clickX      float32
		clickY      float32
		wantBarHit  bool
		wantStillOn bool
	}{
		{"left-click-title-bar", SheetLeft, 700, 10, true, true},
		{"right-click-title-bar", SheetRight, 100, 10, true, true},
		{"left-click-scrim", SheetLeft, 700, 300, false, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			open, barHit := true, false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Fill().Children(func() {
					TitleBar(c, func() {
						ui.Spacer(c)
						if Button(c, "Bar action").Clicked() {
							barHit = true
						}
					})
				})
				Sheet(c, &open, test.side, 300, func() { ui.Text(c, "Sheet body") })
			}, 800, 600)
			tt.Frame()
			body, _ := tt.Find("Sheet body")
			if body.Y < theme.Metrics.HeaderHeight {
				t.Fatalf("sheet starts inside the title bar: %v", body)
			}
			if test.wantBarHit {
				bar, _ := tt.Find("Bar action")
				tt.ClickAt(bar.X+bar.W/2, bar.Y+bar.H/2)
			} else {
				tt.ClickAt(test.clickX, test.clickY)
			}
			if barHit != test.wantBarHit || open != test.wantStillOn {
				t.Fatalf("barHit=%v open=%v", barHit, open)
			}
		})
	}
}
