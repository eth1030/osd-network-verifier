package gcpverifier

import (
	"testing"

	"github.com/openshift/osd-network-verifier/pkg/clients/gcp"
	"github.com/openshift/osd-network-verifier/pkg/mocks"
	gomock "go.uber.org/mock/gomock"
	computev1 "google.golang.org/api/compute/v1"
)

func TestGcpVerifier_findUnreachableEndpoints(t *testing.T) {
	// want to test output with fake serial console outputs
	mockCtrl := gomock.NewController(t)
	// call finish at end to assert mock's expectations
	defer mockCtrl.Finish()
	// obtain mock gcp client
	FakeGCPCli := mocks.NewMockGCPClient(mockCtrl)
	out := &computev1.SerialPortOutput{
		Contents: "fake output",
	}
	FakeGCPCli.EXPECT().GetInstancePorts(gomock.Any(), gomock.Any(), gomock.Any()).Return(out, nil)

	cli := GcpVerifier{
		GcpClient: gcp.Client{
			&computev1.Service,
		},
	}

	err := FakeGCPCli.findUnreachableEndpoints("fake project", "fake zone", "fake instance")
}
