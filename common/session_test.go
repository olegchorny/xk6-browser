package common

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/grafana/xk6-browser/log"
	"github.com/grafana/xk6-browser/storage"
	"github.com/grafana/xk6-browser/tests/ws"

	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/cdp"
	cdppage "github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionCreateSession(t *testing.T) {
	const (
		cdpTargetID         = "target_id_0123456789"
		cdpBrowserContextID = "browser_context_id_0123456789"

		targetAttachedToTargetEvent = `
		{
			"sessionId": "session_id_0123456789",
			"targetInfo": {
				"targetId": "target_id_0123456789",
				"type": "page",
				"title": "",
				"url": "about:blank",
				"attached": true,
				"browserContextId": "browser_context_id_0123456789"
			},
			"waitingForDebugger": false
		}`

		targetAttachedToTargetResult = `
		{
			"sessionId":"session_id_0123456789"
		}
		`
	)

	cmdsReceived := make([]cdproto.MethodType, 0)
	handler := func(conn *websocket.Conn, msg *cdproto.Message, writeCh chan cdproto.Message, done chan struct{}) {
		if msg.SessionID != "" && msg.Method != "" {
			switch msg.Method {
			case cdproto.MethodType(cdproto.CommandPageEnable):
				writeCh <- cdproto.Message{
					ID:        msg.ID,
					SessionID: msg.SessionID,
				}
				close(done) // We're done after receiving the Page.enable command
			}
		} else if msg.Method != "" {
			switch msg.Method {
			case cdproto.MethodType(cdproto.CommandTargetSetDiscoverTargets):
				writeCh <- cdproto.Message{
					ID:        msg.ID,
					SessionID: msg.SessionID,
					Result:    easyjson.RawMessage([]byte("{}")),
				}
			case cdproto.MethodType(cdproto.CommandTargetAttachToTarget):
				writeCh <- cdproto.Message{
					Method: cdproto.EventTargetAttachedToTarget,
					Params: easyjson.RawMessage([]byte(targetAttachedToTargetEvent)),
				}
				writeCh <- cdproto.Message{
					ID:        msg.ID,
					SessionID: msg.SessionID,
					Result:    easyjson.RawMessage([]byte(targetAttachedToTargetResult)),
				}
			}
		}
	}

	server := ws.NewServer(t, ws.WithCDPHandler("/cdp", handler, &cmdsReceived))

	t.Run("send and recv session commands", func(t *testing.T) {
		ctx := context.Background()
		url, _ := url.Parse(server.ServerHTTP.URL)
		wsURL := fmt.Sprintf("ws://%s/cdp", url.Host)
		conn, err := NewConnection(ctx, wsURL, log.NewNullLogger())

		if assert.NoError(t, err) {
			session, err := conn.createSession(&target.Info{
				Type:             "page",
				TargetID:         cdpTargetID,
				BrowserContextID: cdpBrowserContextID,
			})

			if assert.NoError(t, err) {
				action := cdppage.Enable()
				err := action.Do(cdp.WithExecutor(ctx, session))

				require.NoError(t, err)
				require.Equal(t, []cdproto.MethodType{
					cdproto.CommandTargetAttachToTarget,
					cdproto.CommandPageEnable,
				}, cmdsReceived)
			}

			conn.Close()
		}
	})
}

type slowMediaStorer struct {
	t      time.Duration
	doneCh chan struct{}
}

func (s *slowMediaStorer) Store(f storage.MediaFile) error {
	time.Sleep(s.t)
	close(s.doneCh)
	return nil
}

func TestSessionStoreMedia(t *testing.T) { //nolint:tparallel
	t.Parallel()

	// We can not use t.Parallel for each test case
	// as we are overwritting the default "grace period"
	// TO for media storing goroutines after the session
	// is closed.

	t.Run("should wait for media storing goroutines", func(t *testing.T) {
		doneCh := make(chan struct{})
		mediaStorer := &slowMediaStorer{
			t:      50 * time.Millisecond,
			doneCh: doneCh,
		}

		// Overwrite default media storing TO
		mediaStoringTO = 100 * time.Millisecond

		ctx := context.Background()
		sess := NewSession(ctx, nil, "sessID", "targetID", mediaStorer, log.NewNullLogger())

		// Start a "slow" media storing action
		err := sess.StoreMedia("test", nil, "test")
		require.NoError(t, err)

		// Session close will wait for any live uploading
		// goroutines, so we have to call it in a separate
		// goroutine and wait for either it to finish or
		// the slow media storer to finish
		closedCh := make(chan struct{})
		go func() {
			sess.close()
			close(closedCh)
		}()

		select {
		case <-closedCh:
			t.Fatal("session was closed before waiting for media storing")
		case <-doneCh:
		}
	})

	t.Run("should trigger media storing goroutines TO", func(t *testing.T) {
		doneCh := make(chan struct{})
		mediaStorer := &slowMediaStorer{
			t:      100 * time.Millisecond,
			doneCh: doneCh,
		}

		// Overwrite default media storing TO
		mediaStoringTO = 50 * time.Millisecond

		ctx := context.Background()
		sess := NewSession(ctx, nil, "sessID", "targetID", mediaStorer, log.NewNullLogger())

		// Start a "slow" media storing action
		err := sess.StoreMedia("test", nil, "test")
		require.NoError(t, err)

		// Session close will wait for any live uploading
		// goroutines, so we have to call it in a separate
		// goroutine and wait for either it to finish or
		// the slow media storer to finish
		closedCh := make(chan struct{})
		go func() {
			sess.close()
			close(closedCh)
		}()

		select {
		case <-closedCh:
		case <-doneCh:
			t.Fatal("media store unexpectedly finished before TO")
		}
	})
}
