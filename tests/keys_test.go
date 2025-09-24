package containers_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/jsnfwlr/keyper-cli/tests/containers"
)

func noPreStart(ts *containers.TestStack) error {
	return nil
}

func preStartSSHServer(ts *containers.TestStack) error {
	fromTar, _, err := ts.Client.CopyFromContainer(ts.Context, ts.SSHServer.Instance.ID, "/etc/ssh/sshd_config")
	if err != nil {
		fmt.Printf("preStartSSHServer - error 1: %v", err)
		return err
	}

	defer func() {
		_ = fromTar.Close()
	}()

	from := tar.NewReader(fromTar)
	if _, err := from.Next(); err != nil {
		fmt.Printf("preStartSSHServer - error 2: %v", err)
		return err
	}

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, from)
	if err != nil {
		fmt.Printf("preStartSSHServer - error 3: %v", err)
		return err
	}

	fileContent := bytes.ReplaceAll(buf.Bytes(), []byte("#PubkeyAuthentication yes"), []byte("PubkeyAuthentication yes"))
	fileContent = bytes.ReplaceAll(fileContent, []byte("#AuthorizedKeysCommand none"), []byte("AuthorizedKeysCommand /bin/sh /etc/ssh/auth.sh %u %f"))
	fileContent = bytes.ReplaceAll(fileContent, []byte("#AuthorizedKeysCommandUser nobody"), []byte("AuthorizedKeysCommandUser root"))

	buffer, err := containers.TarFile("/etc/ssh/sshd_config", fileContent, 0o0644)
	if err != nil {
		return err
	}

	err = ts.Client.CopyToContainer(ts.Context, ts.SSHServer.Instance.ID, "/", buffer, container.CopyToContainerOptions{CopyUIDGID: false})
	if err != nil {
		fmt.Printf("preStartSSHServer - error 6: %v", err)
		return err
	}

	auth, err := os.ReadFile(filepath.Join(ts.CurrentDir, "testdata", "ssh-server", "auth.sh"))
	if err != nil {
		fmt.Printf("preStartSSHServer - error 7: %v", err)
		return err
	}

	buffer2, err := containers.TarFile("/etc/ssh/auth.sh", auth, 0o0700)
	if err != nil {
		fmt.Printf("preStartSSHServer - error 8: %v", err)
		return err
	}

	err = ts.Client.CopyToContainer(ts.Context, ts.SSHServer.Instance.ID, "/", buffer2, container.CopyToContainerOptions{CopyUIDGID: false})
	if err != nil {
		fmt.Printf("preStartSSHServer - error 9: %v", err)
		return err
	}

	return nil
}

func TestKeys(t *testing.T) {
	// Initialize the test stack
	testStack, err := containers.Start(noPreStart, noPreStart, preStartSSHServer, noPreStart, noPreStart, noPreStart)
	if err != nil {
		t.Fatalf("Failed to create test stack: %v", err)
	}

	// Stop the test stack after use
	defer func() {
		if err := testStack.Stop(); err != nil {
			t.Fatalf("Failed to stop test stack: %v", err)
		}
	}()
}

//
// n := TestContainer{
