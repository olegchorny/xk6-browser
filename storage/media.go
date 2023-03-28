package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	mediaStoreMaxRetries     = 3
	mediaStoreBackoffInitial = 1 * time.Second
	mediaStoreBackoffFactor  = 2
)

// MediaFile represents a media file such as a screenshot.
type MediaFile struct {
	Path        string
	Buf         []byte
	ContentType string
}

// MediaStorer represents a media storing component that can be
// local or remote.
type MediaStorer interface {
	Store(MediaFile) error
}

// NewMediaStorer is a factory method for MediaStorer based on
// the environment variable K6_BROWSER_MEDIA_URL, which determines
// if the MediaStorer should be local or remote.
func NewMediaStorer() MediaStorer {
	url, isRemote := isRemoteMediaStorage()
	if isRemote {
		return &RemoteMediaStorer{
			http: &http.Client{Transport: &http.Transport{}},
			url:  url,
		}
	}

	return &LocalMediaStorer{}
}

// LocalMediaStorer represents a MediaStorer which stores media
// files in the local machine.
type LocalMediaStorer struct{}

// Store stores the media file locally to the given path.
func (l *LocalMediaStorer) Store(f MediaFile) error {
	dir := filepath.Dir(f.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating local media directory %q: %w", dir, err)
	}
	if err := ioutil.WriteFile(f.Path, f.Buf, 0o644); err != nil {
		return fmt.Errorf("saving local media to %q: %w", f.Buf, err)
	}

	return nil
}

// RemoteMediaStorer represents a MediaStorer which pushes media
// files to a remote service through the defined URL.
type RemoteMediaStorer struct {
	http *http.Client
	url  string
}

// Store pushes the media file to the remote store.
func (r *RemoteMediaStorer) Store(f MediaFile) error {
	// Concatenate the file name with the remote store base URL
	url, err := url.JoinPath(r.url, path.Base(f.Path))
	if err != nil {
		return fmt.Errorf("building remote store object URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(f.Buf))
	if err != nil {
		return fmt.Errorf("building remote store PUT request: %w", err)
	}
	req.Header.Set("Content-Type", f.ContentType)

	return retry(func() error {
		resp, err := r.http.Do(req)
		if err != nil {
			return fmt.Errorf("performing PUT request to remote store: %w", err)
		}
		defer resp.Body.Close() //nolint:errcheck
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
		}

		return nil
	})
}

func retry(do func() error) (err error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	backoff := mediaStoreBackoffInitial + time.Duration(r.Int63n(1000))*time.Millisecond

	for i := 0; i < mediaStoreMaxRetries; i++ {
		err = do()
		if err == nil {
			return nil
		}

		time.Sleep(backoff)
		backoff *= time.Duration(mediaStoreBackoffFactor)
	}

	return
}

func isRemoteMediaStorage() (url string, isRemote bool) {
	url, isRemote = os.LookupEnv("K6_BROWSER_MEDIA_URL")
	return
}
