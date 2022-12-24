//go:build docker

package fwa

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

func setupFaktory() (func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() error { return nil }, fmt.Errorf("could not construct pool: %s", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return func() error { return nil }, fmt.Errorf("could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Repository: "contribsys/faktory",
		Name:       "fwa_test",
		Tag:        "latest",
		Cmd:        []string{"/faktory", "-b", ":7419", "-w", ":7420"},
		Env:        []string{},
		PortBindings: map[dc.Port][]dc.PortBinding{
			"7419/tcp": {{HostPort: "7419", HostIP: "127.0.0.1"}},
			"7420/tcp": {{HostPort: "7420", HostIP: "127.0.0.1"}},
		},
	}

	resource, err := pool.RunWithOptions(options, func(config *dc.HostConfig) {
		config.AutoRemove = true
		config.PublishAllPorts = false
	})
	if err != nil {
		return func() error { return nil }, fmt.Errorf("could not start resource: %w", err)
	}

	stopFunc := func() error {
		fmt.Println("cleaning up docker resources")
		return pool.Purge(resource)
	}

	if err := resource.Expire(120); err != nil {
		return stopFunc, fmt.Errorf("could not expire resource: %w", err)
	}

	pool.MaxWait = 10 * time.Second
	if err := pool.Retry(func() error {
		url := fmt.Sprintf("http://localhost:%s/", resource.GetPort("7420/tcp"))
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error retrieving URL: %v\n", err)
			return err
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Status code: %d", resp.StatusCode)
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		return stopFunc, fmt.Errorf("could not connect to docker: %w", err)
	}

	return stopFunc, nil
}

func TestMain(m *testing.M) {
	stop, err := setupFaktory()
	if err != nil {
		panic(err)
	}
	defer stop()
	os.Exit(m.Run())
}
