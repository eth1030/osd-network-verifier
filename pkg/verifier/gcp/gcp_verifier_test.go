package gcpverifier

import (
	"testing"

	"github.com/openshift/osd-network-verifier/pkg/clients/gcp"
	"github.com/openshift/osd-network-verifier/pkg/mocks"
	"github.com/openshift/osd-network-verifier/pkg/probes/curl"
	gomock "go.uber.org/mock/gomock"
	computev1 "google.golang.org/api/compute/v1"
)

func TestGcpVerifier_findUnreachableEndpoints(t *testing.T) {
	/* want to test parsing for findUnreachableEndpoints so want to test multiple fake serial port outputs
	this means we need to create a fake serial port output and a fake gcp client */
	mockCtrl := gomock.NewController(t)
	// call finish at end to assert mock's expectations
	defer mockCtrl.Finish()
	// obtain mock gcp client
	FakeGCPCli := mocks.NewMockGCPClient(mockCtrl)
	// create fake output
	out := &computev1.SerialPortOutput{
		Contents: "fake output",
	}
	// expect GetInstancePorts to be called with any arguments and return fake output
	FakeGCPCli.EXPECT().GetInstancePorts(gomock.Any(), gomock.Any(), gomock.Any()).Return(out, nil)

	cli := GcpVerifier{
		GcpClient: gcp.Client{},
	}

	// need to set client to fake GCP client
	var clients gcp.GCPClient
	clients = FakeGCPCli
	cli.GcpClient = clients

	err := cli.findUnreachableEndpoints("emhammon-test", "us-east-1", "fake instance", curl.Probe{})
	if err != nil {
		t.Errorf("findUnreachableEndpoints() got error %v", err)
	}
	if !cli.Output.IsSuccessful() {
		t.Errorf("Success %v", cli.Output)
	}
}
