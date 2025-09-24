package containers_test

import (
	"testing"
	"time"

	tests "github.com/jsnfwlr/keyper-cli/tests/containers"
)

func noPrep(ts *tests.TestStack) error {
	return nil
}

func TesTheContainersActuallyWork(t *testing.T) {
	// Initialize the test stack

	testStack, err := tests.Start(noPrep, noPrep, noPrep, noPrep, noPrep, noPrep)
	if err != nil {
		t.Fatalf("Failed to create test stack: %v", err)
	}

	// Stop the test stack after use
	defer func() {
		if err := testStack.Stop(); err != nil {
			t.Fatalf("Failed to stop test stack: %v", err)
		}
	}()

	inspect, err := testStack.Client.ContainerInspect(testStack.Context, testStack.SSHServer.Instance.ID)
	if err != nil {
		t.Fatalf("Failed to inspect SSHServer container: %v", err)
	}
	if inspect.State.Running {
		t.Logf("SSHServer Container %s is running", testStack.SSHServer.Instance.ID)
	} else {
		t.Fatalf("SSHServer Container %s is not running", testStack.SSHServer.Instance.ID)
	}

	inspect2, err := testStack.Client.ContainerInspect(testStack.Context, testStack.Keyper.Instance.ID)
	if err != nil {
		t.Fatalf("Failed to inspect Keyper container: %v", err)
	}
	if inspect2.State.Running {
		t.Logf("Keyper Container %s is running", testStack.Keyper.Instance.ID)
	} else {
		t.Fatalf("Keyper Container %s is not running", testStack.Keyper.Instance.ID)
	}

	err = testStack.Stop()
	if err != nil {
		t.Fatalf("Failed to stop test stack: %v", err)
	}
	t.Logf("Test stack stopped successfully")

	time.Sleep(5 * time.Second)

	// Check if the containers are stopped
	_, err = testStack.Client.ContainerInspect(testStack.Context, testStack.SSHServer.Instance.ID)
	if err == nil {
		t.Fatalf("SSHServer Container %s is still running", testStack.SSHServer.Instance.ID)
	} else {
		t.Logf("SSHServer Container %s is stopped", testStack.SSHServer.Instance.ID)
	}

	// Check if the containers are stopped
	_, err = testStack.Client.ContainerInspect(testStack.Context, testStack.Keyper.Instance.ID)
	if err == nil {
		t.Fatalf("Keyper Container %s is still running", testStack.Keyper.Instance.ID)
	} else {
		t.Logf("Keyper Container %s is stopped", testStack.Keyper.Instance.ID)
	}
}

//
// n := TestContainer{
