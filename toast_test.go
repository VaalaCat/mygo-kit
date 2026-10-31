package kit

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// TestToast 验证通知的关闭按钮是 20 见方的紧凑按钮，点击后通知消失。
func TestToast(t *testing.T) {
	// 1. 每种语义各弹出一条通知。
	tests := []struct {
		name string
		kind ToastKind
	}{
		{"info", ToastInfo},
		{"success", ToastSuccess},
		{"error", ToastError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shown := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if !shown {
					shown = true
					ShowToast(c, ToastOpt{Title: "Saved", Description: "detail", Kind: test.kind})
				}
				ToastHost(c)
			}, 600, 400)
			tt.Frame()
			tt.Frame()
			closeButton, ok := tt.Find("Close")
			if !ok || closeButton.W != 20 || closeButton.H != 20 {
				t.Fatalf("close=%v ok=%v texts=%v", closeButton, ok, tt.Texts())
			}
			if err := tt.Click("Close"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if tt.HasText("Saved") {
				t.Fatal("toast still shown")
			}
		})
	}
}

// TestToastPlacement 验证通知在右上角、顶栏下方留 8，最新的在最上面；没有顶栏时从窗口顶部留 8。
func TestToastPlacement(t *testing.T) {
	// 1. 顶栏高度取 36 与系统顶栏的较大值；位置含 1 的边线和 12 的内边距；系统顶栏更高时跟随系统；没有 TitleBar 时不让出顶栏位置。
	tests := []struct {
		name     string
		titleBar bool
		bar      float32
		titles   []string
		wantY    float32
	}{
		{"no-title-bar", false, 0, []string{"only"}, 8 + 1 + 12},
		{"default-title-bar", true, 0, []string{"only"}, 36 + 8 + 1 + 12},
		{"taller-system-title-bar", true, 52, []string{"only"}, 52 + 8 + 1 + 12},
		{"newest-on-top", true, 0, []string{"older", "newer"}, 36 + 8 + 1 + 12},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shown := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if test.titleBar {
					TitleBar(c, func() {})
				}
				if !shown {
					shown = true
					for _, title := range test.titles {
						ShowToast(c, ToastOpt{Title: title, Kind: ToastInfo})
					}
				}
				ToastHost(c)
			}, 600, 400)
			tt.SetTitleBar(ui.TitleBar{Height: test.bar})
			tt.Frame()
			first, ok := tt.Find(test.titles[len(test.titles)-1])
			if !ok {
				t.Fatalf("toast missing: %v", tt.Texts())
			}
			closeButton, _ := tt.Find("Close")
			if closeButton.X+closeButton.W != 600-8-1-12 {
				t.Fatalf("right edge %v", closeButton.X+closeButton.W)
			}
			if len(test.titles) == 1 && closeButton.Y != test.wantY {
				t.Fatalf("top %v want %v", closeButton.Y, test.wantY)
			}
			if len(test.titles) == 2 {
				older, _ := tt.Find(test.titles[0])
				if !(first.Y < older.Y) {
					t.Fatalf("newer %v is not above older %v", first.Y, older.Y)
				}
			}
		})
	}
}

// TestToastNeedsNoManualRedraw 验证空闲窗口里添加的通知立即出现，到期后自己消失。
func TestToastNeedsNoManualRedraw(t *testing.T) {
	// 1. ToastHost 先绘制，通知是在它之后添加的，只有请求下一帧才能显示。
	tests := []struct {
		name    string
		kind    ToastKind
		advance time.Duration
	}{
		{"info-expires-after-4s", ToastInfo, 5 * time.Second},
		{"error-stays-longer-and-expires-after-8s", ToastError, 9 * time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			add := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				ToastHost(c)
				if add {
					add = false
					ShowToast(c, ToastOpt{Title: "Saved", Kind: test.kind})
				}
			}, 600, 400)
			add = true
			tt.Frame()
			if !tt.HasText("Saved") {
				t.Fatalf("toast not shown without another frame: %v", tt.Texts())
			}
			if test.kind == ToastError {
				tt.Advance(5 * time.Second)
				if !tt.HasText("Saved") {
					t.Fatal("error toast expired too early")
				}
			}
			tt.Advance(test.advance)
			if tt.HasText("Saved") {
				t.Fatal("toast did not expire in an idle window")
			}
		})
	}
}
