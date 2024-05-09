package postal

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/reiver/go-telnet"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	smtpResource, smtpPool := mailHogUp()

	// run tests
	code := m.Run()

	// clean up
	if err := smtpPool.Purge(smtpResource); err != nil {
		log.Fatalf("could not purge smtp resource: %s", err)
	}

	os.Exit(code)
}

func mailHogUp() (*dockertest.Resource, *dockertest.Pool) {
	// connect to docker; fail if docker not running
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository:   "jcalonso/mailhog",
		Tag:          "latest",
		ExposedPorts: []string{"1025"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1025": {
				{HostIP: "0.0.0.0", HostPort: "1026"},
			},
		},
	}

	// get a resource (docker image)
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var caller = telnet.StandardCaller
		retryErr := telnet.DialToAndCall("localhost:1026", caller)
		if retryErr != nil {
			log.Println("Error:", retryErr)
			return retryErr
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to mailserver: %s", err)
	}

	return resource, pool
}
