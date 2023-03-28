package storage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoteMediaStorerStore(t *testing.T) {
	t.Parallel()

	// verifyReq contains common verifications
	// for the performed request to test server
	verifyReq := func(r *http.Request) {
		require.Equal(t, "/object.jpeg", r.URL.Path)
		require.Equal(t, "image/jpeg", r.Header.Get("Content-Type"))
	}

	// nReq is used to return
	// an error from 'failing handler'
	// test case only on first request
	var nReq int

	testCases := []struct {
		name       string
		srvHandler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "happy path",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				verifyReq(r)
			},
		},
		{
			name: "should retry",
			srvHandler: func(w http.ResponseWriter, r *http.Request) {
				verifyReq(r)
				if nReq == 0 {
					w.WriteHeader(http.StatusInternalServerError)
					nReq++
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(http.HandlerFunc(tc.srvHandler))
			defer ts.Close()

			storer := &RemoteMediaStorer{
				http: &http.Client{Transport: &http.Transport{}},
				url:  ts.URL,
			}

			f := MediaFile{
				Path:        "/path/to/object.jpeg",
				Buf:         []byte("someData"),
				ContentType: "image/jpeg",
			}

			require.NoError(t, storer.Store(f))
		})
	}
}
