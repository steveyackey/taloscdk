package taloscdk

import (
	"os"
	"strings"

	"github.com/aws/jsii-runtime-go"
)

// LoadAndTransformMachineConfig takes a Talos cluster config file and replaces the endpoint with the
// correct hostname or IP based on what CDK generates.
// To get started, you can try running `talosctl gen config talos https://talos.cluster:6443`
// and use the controlplane.yaml as the filename
func LoadConfig(fileName string) (*string, error) {
	config, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return jsii.String(string(config)), nil
}

// TransformConfig replaces an initialEndpoint with a replacementEndpoint
func TransformConfig(config *string, initialEndpoint string, replacementEndpoint string) *string {
	return jsii.String(strings.ReplaceAll(*config, initialEndpoint, replacementEndpoint))
}
