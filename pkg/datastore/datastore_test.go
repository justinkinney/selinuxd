package datastore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func getNewStorePath(t *testing.T) (dspath string, cleanup func()) {
	d, err := ioutil.TempDir("", "store")
	if err != nil {
		t.Fatalf("Couldn't create tmpfile")
	}
	return filepath.Join(d, "policy.db"), func() {
		os.RemoveAll(d)
	}
}

func getNewStore(path string, t *testing.T) (datastore DataStore, cleanup func()) {
	ds, err := New(path)
	if err != nil {
		t.Fatalf("Couldn't create tmpfile")
	}
	return ds, func() {
		ds.Close()
	}
}

func TestDataStore(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		generateTmpDir bool
		wantErr        bool
	}{
		{"Basic usage", "", true, false},
		{"unexistent dir", "/unexistent-dir/unexistent-path", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.generateTmpDir {
				d, err := ioutil.TempDir("", "store")
				if err != nil {
					t.Errorf("Couldn't create tmpfile")
					return
				}
				defer os.RemoveAll(d)
				path = filepath.Join(d, "policy.db")
			} else {
				path = tt.path
			}

			got, err := New(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("New() DataStore store object should only be nil if an error happened, wantErr %v", tt.wantErr)
				return
			}
		})
	}
}

func TestStatusProbe(t *testing.T) {
	type Args struct {
		policy string
		status StatusType
		msg    string
	}
	args := Args{"my-policy", "installed", "all is good"}

	path, filecleanup := getNewStorePath(t)
	defer filecleanup()
	ds, dscleanup := getNewStore(path, t)
	defer dscleanup()

	// We should be able to write a status correctly
	if err := ds.PutStatus(args.policy, args.status, args.msg); err != nil {
		t.Errorf("DataStore.PutStatus() error = %v", err)
	}

	// We should be able to read a status correctly
	status, msg, err := ds.GetStatus(args.policy)
	if err != nil {
		t.Errorf("DataStore.GetStatus() error = %v", err)
	}
	if args.status != status {
		t.Errorf("DataStore.GetStatus() status didn't match. got: %s, expected: %s", status, args.status)
	}
	if args.msg != msg {
		t.Errorf("DataStore.GetStatus() msg didn't match. got: %s, expected: %s", msg, args.msg)
	}
}

func TestStatusProbeReadOnly(t *testing.T) {
	type Args struct {
		policy string
		status StatusType
		msg    string
	}
	args := Args{"my-policy", "installed", "all is good"}

	path, filecleanup := getNewStorePath(t)
	defer filecleanup()
	ds, dscleanup := getNewStore(path, t)
	defer dscleanup()
	rods := ds.GetReadOnly()

	// We should be able to write a status correctly
	if err := ds.PutStatus(args.policy, args.status, args.msg); err != nil {
		t.Errorf("DataStore.PutStatus() error = %v", err)
	}

	// We should be able to read a status correctly with the read-only interface
	status, msg, err := rods.GetStatus(args.policy)
	if err != nil {
		t.Errorf("DataStore.GetStatus() error = %v", err)
	}
	if args.status != status {
		t.Errorf("DataStore.GetStatus() status didn't match. got: %s, expected: %s", status, args.status)
	}
	if args.msg != msg {
		t.Errorf("DataStore.GetStatus() msg didn't match. got: %s, expected: %s", msg, args.msg)
	}
}

func TestListPolicies(t *testing.T) {
	type Args struct {
		policy string
		status StatusType
		msg    string
	}
	argsList := []Args{
		{"my-policy-1", "installed", "all is good"},
		{"my-policy-2", "installed", "all is good"},
		{"my-policy-3", "installed", "all is good"},
	}

	path, filecleanup := getNewStorePath(t)
	defer filecleanup()
	ds, dscleanup := getNewStore(path, t)
	defer dscleanup()

	for _, args := range argsList {
		// We should be able to write a status correctly
		if err := ds.PutStatus(args.policy, args.status, args.msg); err != nil {
			t.Errorf("DataStore.PutStatus() error = %v", err)
		}
	}

	policies, err := ds.List()
	if err != nil {
		t.Errorf("DataStore.List() error = %v", err)
	}
	if len(policies) != len(argsList) {
		t.Errorf("DataStore.List() didn't output the expected number of policies. Got %d, Expected %d",
			len(policies), len(argsList))
	}
}

func TestRemovePolicy(t *testing.T) {
	type Args struct {
		policy string
		status StatusType
		msg    string
	}
	args := Args{"my-policy-1", "installed", "all is good"}

	path, filecleanup := getNewStorePath(t)
	defer filecleanup()
	ds, dscleanup := getNewStore(path, t)
	defer dscleanup()

	// We should be able to write a status correctly
	if err := ds.PutStatus(args.policy, args.status, args.msg); err != nil {
		t.Errorf("DataStore.PutStatus() error = %v", err)
	}

	policies, err := ds.List()
	if err != nil {
		t.Errorf("DataStore.List() error = %v", err)
	}
	if len(policies) != 1 {
		t.Errorf("DataStore.List() didn't output the expected number of policies. Got %d, Expected %d",
			len(policies), 1)
	}

	// Remove the policy
	if err := ds.Remove(args.policy); err != nil {
		t.Errorf("DataStore.Remove() error = %v", err)
	}

	policies, err = ds.List()
	if err != nil {
		t.Errorf("DataStore.List() error = %v", err)
	}
	if len(policies) != 0 {
		t.Errorf("DataStore.List() didn't output the expected number of policies. Got %d, Expected %d",
			len(policies), 0)
	}
}
