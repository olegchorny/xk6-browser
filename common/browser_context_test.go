package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/xk6-browser/common/js"
	"github.com/grafana/xk6-browser/k6ext"
	"github.com/grafana/xk6-browser/k6ext/k6test"
	"github.com/grafana/xk6-browser/log"
)

func TestNewBrowserContext(t *testing.T) {
	t.Run("add_web_vital_js_scripts_to_context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		logger := log.NewNullLogger()
		b := newBrowser(ctx, cancel, nil, NewLaunchOptions(), logger)

		vu := k6test.NewVU(t)
		ctx = k6ext.WithVU(ctx, vu)

		bc, err := NewBrowserContext(ctx, b, "some-id", nil, nil)
		require.NoError(t, err)

		webVitalIIFEScriptFound := false
		webVitalInitScriptFound := false
		for _, script := range bc.evaluateOnNewDocumentSources {
			switch script {
			case js.WebVitalIIFEScript:
				webVitalIIFEScriptFound = true
			case js.WebVitalInitScript:
				webVitalInitScriptFound = true
			default:
				assert.Fail(t, "script is neither WebVitalIIFEScript nor WebVitalInitScript")
			}
		}

		assert.True(t, webVitalIIFEScriptFound, "WebVitalIIFEScript was not initialized in the context")
		assert.True(t, webVitalInitScriptFound, "WebVitalInitScript was not initialized in the context")
	})
}
