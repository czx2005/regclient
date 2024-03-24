package ocidir

import (
	"context"
	"errors"
	"testing"

	"github.com/regclient/regclient/internal/rwfs"
	"github.com/regclient/regclient/types/errs"
	"github.com/regclient/regclient/types/ref"
)

func TestTag(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fsOS := rwfs.OSNew("")
	fsMem := rwfs.MemNew()
	err := rwfs.MkdirAll(fsMem, "testdata/regctl", 0777)
	if err != nil {
		t.Fatalf("failed to setup memfs dir: %v", err)
	}
	err = rwfs.CopyRecursive(fsOS, "testdata/regctl", fsMem, "testdata/regctl")
	if err != nil {
		t.Fatalf("failed to setup memfs copy: %v", err)
	}
	oMem := New(WithFS(fsMem))
	tRef := "ocidir://testdata/regctl"
	r, err := ref.New(tRef)
	if err != nil {
		t.Fatalf("failed to parse ref %s: %v", tRef, err)
	}
	rCp := r

	t.Run("TagList", func(t *testing.T) {
		exTags := []string{"broken", "latest", "v0.3", "v0.3.10"}
		tl, err := oMem.TagList(ctx, r)
		if err != nil {
			t.Fatalf("failed to retrieve tag list: %v", err)
		}
		tlTags, err := tl.GetTags()
		if err != nil {
			t.Fatalf("failed to get tags: %v", err)
		}
		if !cmpSliceString(exTags, tlTags) {
			t.Errorf("unexpected tag list, expected %v, received %v", exTags, tlTags)
		}
	})

	t.Run("TagDelete", func(t *testing.T) {
		exTags := []string{"broken", "v0.3"}
		rCp.Tag = "missing"
		err := oMem.TagDelete(ctx, rCp)
		if err == nil || !errors.Is(err, errs.ErrNotFound) {
			t.Errorf("deleting missing tag %s: %v", rCp.CommonName(), err)
		}
		rCp.Tag = "latest"
		err = oMem.TagDelete(ctx, rCp)
		if err != nil {
			t.Errorf("failed to delete tag %s: %v", rCp.CommonName(), err)
		}
		rCp.Tag = "v0.3.10"
		err = oMem.TagDelete(ctx, rCp)
		if err != nil {
			t.Errorf("failed to delete tag %s: %v", rCp.CommonName(), err)
		}

		tl, err := oMem.TagList(ctx, r)
		if err != nil {
			t.Fatalf("failed to retrieve tag list: %v", err)
		}
		tlTags, err := tl.GetTags()
		if err != nil {
			t.Errorf("failed to get tags: %v", err)
		}
		if !cmpSliceString(exTags, tlTags) {
			t.Errorf("unexpected tag list, expected %v, received %v", exTags, tlTags)
		}
	})
}
