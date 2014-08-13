package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWatcherNoticesCreateFile(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	watcher, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	f, err := ioutil.TempFile(tmpdir, "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Creating file did not generate event")
	}
}

func TestWatcherNoticesWriteFile(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	f, err := ioutil.TempFile(tmpdir, "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	watcher, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString("asdffdsakjkdsklsdakjaskjadklaskjlsd")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Writing to file did not generate event")
	}
}

func TestWatcherNoticesRemoveFile(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	f, err := ioutil.TempFile(tmpdir, "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	watcher, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Deleting file did not generate event")
	}
}

func TestWatcherNoticesCreatedSubdirectoryAndChangesWithinIt(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	watcher, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(tmpdir)
	if err != nil {
		t.Fatal(err)
	}

	subdir, err := ioutil.TempDir(tmpdir, "watcher_")
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Creating subdirectory did not trigger event")
	}

	f, err := ioutil.TempFile(subdir, "watcher_")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Creating file in subdirectory did not trigger event")
	}
}