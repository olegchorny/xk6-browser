package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/xk6-browser/browser"
	"github.com/grafana/xk6-browser/chromium"
	"github.com/grafana/xk6-browser/k6ext/k6test"
)

func TestBrowserTypeConnect(t *testing.T) {
	// Start a test browser so we can get its WS URL
	// and use it to connect through BrowserType.Connect.
	tb := newTestBrowser(t)
	vu := k6test.NewVU(t)
	bt := chromium.NewBrowserType(vu)
	vu.MoveToVUContext()

	b := bt.Connect(tb.wsURL, nil)
	b.NewPage(nil)
}

func TestBrowserTypeLaunchToConnect(t *testing.T) {
	vu := k6test.NewVU(t)
	tb := newTestBrowser(t)
	bp := newTestBrowserProxy(t, tb)

	// Export WS URL env var
	// pointing to test browser proxy
	t.Setenv("K6_BROWSER_WS_URL", bp.wsURL())

	// We have to call launch method through JS API in Goja
	// to take mapping layer into account, instead of calling
	// BrowserType.Launch method directly
	rt := vu.Runtime()
	root := browser.New()
	mod := root.NewModuleInstance(vu)
	jsMod, ok := mod.Exports().Default.(*browser.JSModule)
	require.Truef(t, ok, "unexpected default mod export type %T", mod.Exports().Default)

	vu.MoveToVUContext()

	require.NoError(t, rt.Set("chromium", jsMod.Chromium))
	_, err := rt.RunString(`
		const b = chromium.launch();
		b.close();
	`)
	require.NoError(t, err)

	// Verify the proxy, which's WS URL was set as
	// K6_BROWSER_WS_URL, has received a connection req
	require.True(t, bp.connected)
	// Verify that no new process pids have been added
	// to pid registry
	require.Len(t, root.PidRegistry.Pids(), 0)
}
