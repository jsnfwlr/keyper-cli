package containers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/jsnfwlr/keyper-cli/tests/containers/console"
)

type TestStack struct {
	Keyper     *TestContainer
	SSHServer  *TestContainer
	SSHClient  *TestContainer
	Network    *TestNetwork
	Client     *client.Client
	Context    context.Context
	CurrentDir string
}

type TestContainer struct {
	Label         string
	Config        container.Config
	HostConfig    container.HostConfig
	NetworkConfig network.NetworkingConfig
	Instance      container.CreateResponse
}

type TestNetwork struct {
	Name     string
	Instance network.CreateResponse
	Config   network.NetworkingConfig
}

type PrepFunc func(testStack *TestStack) (fault error)

func New() (*TestStack, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ts := &TestStack{
		Client:     dc,
		Context:    context.Background(),
		CurrentDir: cwd,
	}

	return ts, nil
}

func Start(kepyerPre, kepyerPost, sshServerPre, sshServerPost, sshClientPre, sshClientPost PrepFunc) (testStack *TestStack, fault error) {
	ts, err := New()
	if err != nil {
		return nil, err
	}

	err = ts.StartNetwork()
	if err != nil {
		return nil, err
	}

	err = ts.StartKeyper(kepyerPre, kepyerPost)
	if err != nil {
		_ = ts.Stop()
		return nil, err
	}

	err = ts.StartSSHServer(sshServerPre, sshServerPost)
	if err != nil {
		_ = ts.Stop()
		return nil, err
	}

	// err = ts.StartSSHClient(sshClientPre, sshClientPost)
	// if err != nil {
	// 	_ = ts.Stop()
	// 	return nil, err
	// }

	return ts, nil
}

func (ts *TestStack) StartNetwork() (fault error) {
	var err error

	label := "keyper-cli-test"

	n := TestNetwork{
		Name: label,
	}

	n.Instance, err = ts.Client.NetworkCreate(ts.Context, n.Name, network.CreateOptions{Driver: "bridge"})
	if err != nil {
		return err
	}

	n.Config = network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			label: {
				NetworkID: n.Instance.ID,
			},
		},
	}

	ts.Network = &n

	return nil
}

func (ts *TestStack) StartKeyper(preStart, postStart PrepFunc) (fault error) {
	label := "keyper-cli-test-keyper"
	img := "dbsentry/keyper:latest"

	err := ts.GetImage(img)
	if err != nil {
		return err
	}

	c := TestContainer{
		Label: label,

		Config: container.Config{
			Hostname: label,
			Image:    img,
			Env: []string{
				fmt.Sprintf("HOSTNAME=%s", label),
				"LDAP_ORGANIZATION_NAME=jsnfwlr",
				"LDAP_DOMAIN=jsnfwlr.com",
				"LDAP_ADMIN_PASSWORD=keyper-cli-test-ldap",
				"FLASK_CONFIG=prod",
				"NGINX_UID=1000",
				"NGINX_GID=1000",
				"CONTAINER_LOG_LEVEL=4",
				"ENVIRONMENT_LOG_LEVEL_KEY=debug",
			},
		},

		HostConfig: container.HostConfig{
			AutoRemove: true,
		},

		NetworkConfig: ts.Network.Config,
	}

	c.Instance, err = ts.Client.ContainerCreate(ts.Context, &c.Config, &c.HostConfig, &c.NetworkConfig, nil, c.Label)
	if err != nil {
		return err
	}

	ts.Keyper = &c

	err = preStart(ts)
	if err != nil {
		return err
	}

	err = ts.Client.ContainerStart(ts.Context, ts.Keyper.Instance.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	err = postStart(ts)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TestStack) StartSSHServer(preStart, postStart PrepFunc) (fault error) {
	label := "keyper-cli-test-ssh-server"
	img := "lscr.io/linuxserver/openssh-server:latest"

	err := ts.GetImage(img)
	if err != nil {
		return err
	}

	c := TestContainer{
		Label: label,

		Config: container.Config{
			Hostname: label,
			Image:    img,
			Env: []string{
				"PUID=1000",
				"PGID=1000",
				"TZ=Etc/UTC",
				"PUBLIC_KEY_URL=",
				"SUDO_ACCESS=false",
				"PASSWORD_ACCESS=false",
				"USER_NAME=keyper-cli-test-ssh-server",
				"USER_PASSWORD=keyper-cli-test-password",
			},
		},

		HostConfig: container.HostConfig{
			AutoRemove: true,
		},

		NetworkConfig: ts.Network.Config,
	}

	c.Instance, err = ts.Client.ContainerCreate(ts.Context, &c.Config, &c.HostConfig, &c.NetworkConfig, nil, c.Label)
	if err != nil {
		return err
	}

	ts.SSHServer = &c

	err = preStart(ts)
	if err != nil {
		return err
	}

	err = ts.Client.ContainerStart(ts.Context, ts.SSHServer.Instance.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	err = postStart(ts)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TestStack) StartSSHClient(preStart, postStart PrepFunc) (fault error) {
	label := "keyper-cli-test-ssh-client"
	img := ""

	err := ts.GetImage(img)
	if err != nil {
		return err
	}

	c := TestContainer{
		Label: label,

		Config: container.Config{
			Hostname: label,
			Image:    img,
			Env:      []string{},
		},

		HostConfig: container.HostConfig{
			AutoRemove: true,
		},

		NetworkConfig: ts.Network.Config,
	}

	c.Instance, err = ts.Client.ContainerCreate(ts.Context, &c.Config, &c.HostConfig, &c.NetworkConfig, nil, c.Label)
	if err != nil {
		return err
	}

	ts.SSHClient = &c

	err = preStart(ts)
	if err != nil {
		return err
	}

	err = ts.Client.ContainerStart(ts.Context, ts.SSHClient.Instance.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	err = postStart(ts)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TestStack) Stop() error {
	timeout := 5

	contStopOpt := container.StopOptions{
		Signal:  "SIGKILL",
		Timeout: &timeout,
	}

	if ts.Keyper != nil {
		if err := ts.Client.ContainerStop(ts.Context, ts.Keyper.Instance.ID, contStopOpt); err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}
		}
	}

	if ts.SSHServer != nil {
		if err := ts.Client.ContainerStop(ts.Context, ts.SSHServer.Instance.ID, contStopOpt); err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}
		}
	}

	if ts.SSHClient != nil {
		if err := ts.Client.ContainerStop(ts.Context, ts.SSHClient.Instance.ID, contStopOpt); err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}
		}
	}

	if ts.Network != nil {
		if err := ts.Client.NetworkRemove(ts.Context, ts.Network.Instance.ID); err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}

func (ts *TestStack) GetImage(imageName string) error {
	found := false

	images, err := ts.Client.ImageList(ts.Context, image.ListOptions{})
	if err != nil {
		return err
	}

	// if the image is tagged with something other than ":latest", we check if the image is already present
	// otherwise skip the check and just pull the image
	if strings.Contains(imageName, ":") && !strings.Contains(imageName, ":latest") {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if tag == imageName {
					fmt.Printf("Found %s - skipping pull\n", tag)
					found = true
					break
				}
			}
		}
	}

	if !found {
		output, err := ts.Client.ImagePull(
			ts.Context,
			imageName,
			image.PullOptions{
				Platform: os.Getenv("GOARCH"),
			},
		)
		if err != nil {
			return fmt.Errorf("error pulling image: %w", err)
		}

		defer func() {
			_ = output.Close()
		}()

		err = console.Print(output)
		if err != nil {
			return fmt.Errorf("error printing image pull output: %w", err)
		}
	}

	return nil
}
