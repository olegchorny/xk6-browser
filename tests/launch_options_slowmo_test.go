package tests

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/xk6-browser/api"
	"github.com/grafana/xk6-browser/common"
)

func TestLaunchOptionsSlowMo(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	t.Run("Page", func(t *testing.T) {
		t.Parallel()
		t.Run("check", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Check(".check", nil)
			})
		})
		t.Run("click", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				err := p.Click("button", nil)
				assert.NoError(t, err)
			})
		})
		t.Run("dblClick", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Dblclick("button", nil)
			})
		})
		t.Run("dispatchEvent", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.DispatchEvent("button", "click", goja.Null(), nil)
			})
		})
		t.Run("emulateMedia", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.EmulateMedia(tb.toGojaValue(struct {
					Media string `js:"media"`
				}{
					Media: "print",
				}))
			})
		})
		t.Run("evaluate", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Evaluate(tb.toGojaValue("() => void 0"))
			})
		})
		t.Run("evaluateHandle", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.EvaluateHandle(tb.toGojaValue("() => window"))
			})
		})
		t.Run("fill", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Fill(".fill", "foo", nil)
			})
		})
		t.Run("focus", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Focus("button", nil)
			})
		})
		t.Run("goto", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				_, err := p.Goto("about:blank", nil)
				require.NoError(t, err)
			})
		})
		t.Run("hover", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Hover("button", nil)
			})
		})
		t.Run("press", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Press("button", "Enter", nil)
			})
		})
		t.Run("reload", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Reload(nil)
			})
		})
		t.Run("setContent", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.SetContent("hello world", nil)
			})
		})
		/*t.Run("setInputFiles", func(t *testing.T) {
			testPageSlowMoImpl(t, tb, func(_ *Browser, p api.Page) {
				p.SetInputFiles(".file", nil, nil)
			})
		})*/
		t.Run("selectOption", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.SelectOption("select", tb.toGojaValue("foo"), nil)
			})
		})
		t.Run("setViewportSize", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.SetViewportSize(nil)
			})
		})
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Type(".fill", "a", nil)
			})
		})
		t.Run("uncheck", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testPageSlowMoImpl(t, tb, func(_ *testBrowser, p api.Page) {
				p.Uncheck(".uncheck", nil)
			})
		})
	})

	t.Run("Frame", func(t *testing.T) {
		t.Parallel()
		t.Run("check", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Check(".check", nil)
			})
		})
		t.Run("click", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				err := f.Click("button", nil)
				assert.NoError(t, err)
			})
		})
		t.Run("dblClick", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Dblclick("button", nil)
			})
		})
		t.Run("dispatchEvent", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.DispatchEvent("button", "click", goja.Null(), nil)
			})
		})
		t.Run("evaluate", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Evaluate(tb.toGojaValue("() => void 0"))
			})
		})
		t.Run("evaluateHandle", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.EvaluateHandle(tb.toGojaValue("() => window"))
			})
		})
		t.Run("fill", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Fill(".fill", "foo", nil)
			})
		})
		t.Run("focus", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Focus("button", nil)
			})
		})
		t.Run("goto", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Goto("about:blank", nil)
			})
		})
		t.Run("hover", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Hover("button", nil)
			})
		})
		t.Run("press", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Press("button", "Enter", nil)
			})
		})
		t.Run("setContent", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.SetContent("hello world", nil)
			})
		})
		/*t.Run("setInputFiles", func(t *testing.T) {
			testFrameSlowMoImpl(t, tb, func(_ *Browser, f api.Frame) {
				f.SetInputFiles(".file", nil, nil)
			})
		})*/
		t.Run("selectOption", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.SelectOption("select", tb.toGojaValue("foo"), nil)
			})
		})
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Type(".fill", "a", nil)
			})
		})
		t.Run("uncheck", func(t *testing.T) {
			t.Parallel()
			tb := newTestBrowser(t, withFileServer())
			testFrameSlowMoImpl(t, tb, func(_ *testBrowser, f api.Frame) {
				f.Uncheck(".uncheck", nil)
			})
		})
	})

	// TODO implement this
	t.Run("ElementHandle", func(t *testing.T) {
	})
}

func testSlowMoImpl(t *testing.T, tb *testBrowser, fn func(*testBrowser)) {
	hooks := common.GetHooks(tb.ctx)
	currentHook := hooks.Get(common.HookApplySlowMo)
	chCalled := make(chan bool, 1)
	defer hooks.Register(common.HookApplySlowMo, currentHook)
	hooks.Register(common.HookApplySlowMo, func(ctx context.Context) {
		currentHook(ctx)
		chCalled <- true
	})

	didSlowMo := false
	go fn(tb)
	select {
	case <-tb.ctx.Done():
	case <-chCalled:
		didSlowMo = true
	}

	require.True(t, didSlowMo, "expected action to have been slowed down")
}

func testPageSlowMoImpl(t *testing.T, tb *testBrowser, fn func(*testBrowser, api.Page)) {
	p := tb.NewPage(nil)

	p.SetContent(`
		<button>a</button>
		<input type="checkbox" class="check">
		<input type="checkbox" checked=true class="uncheck">
		<input class="fill">
		<select>
		<option>foo</option>
		</select>
		<input type="file" class="file">
    	`, nil)
	testSlowMoImpl(t, tb, func(tb *testBrowser) { fn(tb, p) })
}

func testFrameSlowMoImpl(t *testing.T, tb *testBrowser, fn func(bt *testBrowser, f api.Frame)) {
	p := tb.NewPage(nil)

	f := tb.attachFrame(p, "frame1", tb.staticURL("empty.html"))
	f.SetContent(`
		<button>a</button>
		<input type="checkbox" class="check">
		<input type="checkbox" checked=true class="uncheck">
		<input class="fill">
		<select>
		  <option>foo</option>
		</select>
		<input type="file" class="file">
    	`, nil)
	testSlowMoImpl(t, tb, func(tb *testBrowser) { fn(tb, f) })
}
