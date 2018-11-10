package k8s

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDeleteConfigContext(t *testing.T) {
	t.Parallel()

	path := storeConfigToTempFile(t, BASIC_CONFIG_WITH_EXTRA_CONTEXT)
	defer os.Remove(path)

	if err := DeleteConfigContextWithPathE(t, path, "extra_minikube"); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	storedConfig := string(data)
	assert.Equal(t, storedConfig, BASIC_CONFIG)
}

func TestDeleteConfigContextWithAnotherContextRemaining(t *testing.T) {
	t.Parallel()

	path := storeConfigToTempFile(t, BASIC_CONFIG_WITH_EXTRA_CONTEXT_NO_GARBAGE)
	defer os.Remove(path)

	if err := DeleteConfigContextWithPathE(t, path, "extra_minikube"); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	storedConfig := string(data)
	assert.Equal(t, storedConfig, EXPECTED_CONFIG_AFTER_EXTRA_MINIKUBE_DELETED_NO_GARBAGE)
}

func TestRemoveOrphanedClusterAndAuthInfoConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			"TestExtraClusterRemoveOrphanedClusterAndAuthInfoed",
			BASIC_CONFIG_WITH_EXTRA_CLUSTER,
			BASIC_CONFIG,
		},
		{
			"TestExtraAuthInfoRemoveOrphanedClusterAndAuthInfoed",
			BASIC_CONFIG_WITH_EXTRA_AUTH_INFO,
			BASIC_CONFIG,
		},
	}
	for _, testCase := range testCases {
		// Capture range variable to scope within range
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			removeOrphanedClusterAndAuthInfoConfigTestFunc(t, testCase.in, testCase.out)
		})
	}
}

func removeOrphanedClusterAndAuthInfoConfigTestFunc(t *testing.T, inputConfig string, expectedOutputConfig string) {
	path := storeConfigToTempFile(t, inputConfig)
	defer os.Remove(path)

	config := LoadConfigFromPath(path)
	rawConfig, err := config.RawConfig()
	if err != nil {
		t.Fatal(err)
	}
	RemoveOrphanedClusterAndAuthInfoConfig(&rawConfig)
	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false); err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	storedConfig := string(data)
	assert.Equal(t, storedConfig, expectedOutputConfig)
}

func storeConfigToTempFile(t *testing.T, configData string) string {
	escapedTestName := url.PathEscape(t.Name())
	tmpfile, err := ioutil.TempFile("", escapedTestName)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.WriteString(configData); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

// Various example configs used in testing the config manipulation functions

const BASIC_CONFIG = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`

const BASIC_CONFIG_WITH_EXTRA_CLUSTER = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`

const BASIC_CONFIG_WITH_EXTRA_AUTH_INFO = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const BASIC_CONFIG_WITH_EXTRA_CONTEXT = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: extra_minikube
  name: extra_minikube
current-context: extra_minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const BASIC_CONFIG_WITH_EXTRA_CONTEXT_NO_GARBAGE = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: extra_minikube
  name: extra_minikube
- context:
    cluster: extra_minikube
    user: minikube
  name: other_minikube

current-context: extra_minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const EXPECTED_CONFIG_AFTER_EXTRA_MINIKUBE_DELETED_NO_GARBAGE = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: minikube
  name: other_minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`
