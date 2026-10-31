package teldrive

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rclone/rclone/backend/teldrive/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newDirMtimeTestObject records the IDs PATCHed with a modTime and answers
// everything else as a successful upload completion or folder creation.
func newDirMtimeTestObject(t *testing.T, patchStatus int) (*Object, *[]string) {
	t.Helper()
	patched := []string{}
	o := newUploadTestObject(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPatch {
			var update api.UpdateFileInformation
			require.NoError(t, json.NewDecoder(r.Body).Decode(&update))
			require.NotNil(t, update.ModTime)
			assert.Empty(t, update.Name)
			assert.WithinDuration(t, time.Now(), *update.ModTime, 5*time.Second)
			patched = append(patched, r.URL.Path)
			w.WriteHeader(patchStatus)
			_, _ = io.WriteString(w, `{}`)
			return
		}
		_, _ = io.WriteString(w, `{"id":"new","name":"file.txt","kind":"file"}`)
	})
	return o, &patched
}

func TestCompleteUploadTouchesParent(t *testing.T) {
	o, patched := newDirMtimeTestObject(t, http.StatusOK)
	_, err := o.completeUpload(context.Background(), &uploadInfo{uploadID: "upload", dir: "parent"})
	require.NoError(t, err)
	assert.Equal(t, []string{"/api/v1/files/parent"}, *patched)
}

func TestCreateDirTouchesParent(t *testing.T) {
	o, patched := newDirMtimeTestObject(t, http.StatusOK)
	id, err := o.fs.CreateDir(context.Background(), "parent", "child")
	require.NoError(t, err)
	assert.Equal(t, "new", id)
	assert.Equal(t, []string{"/api/v1/files/parent"}, *patched)
}

func TestTouchDirSkipsRoot(t *testing.T) {
	o, patched := newDirMtimeTestObject(t, http.StatusOK)
	_, err := o.completeUpload(context.Background(), &uploadInfo{uploadID: "upload", dir: rootID})
	require.NoError(t, err)
	assert.Empty(t, *patched)
}

func TestTouchDirFailureDoesNotFailUpload(t *testing.T) {
	o, patched := newDirMtimeTestObject(t, http.StatusNotFound)
	file, err := o.completeUpload(context.Background(), &uploadInfo{uploadID: "upload", dir: "parent"})
	require.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, []string{"/api/v1/files/parent"}, *patched)
}
