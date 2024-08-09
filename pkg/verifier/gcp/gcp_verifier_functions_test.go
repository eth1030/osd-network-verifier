package gcpverifier

import (
	"testing"

	"github.com/openshift/osd-network-verifier/pkg/probes"
	"github.com/openshift/osd-network-verifier/pkg/probes/curl"
)

func TestGetTokens(t *testing.T) {
	type args struct {
		consoleOutput string
		probe         probes.Probe
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "tokens in order",
			args: args{
				consoleOutput: "otherinfoNV_CURLJSON_BEGIN\nhello world\nNV_CURLJSON_END\njj",
				probe:         curl.Probe{},
			},
			want: true,
		},
		{
			name: "only start token",
			args: args{
				consoleOutput: "NV_CURLJSON_BEGIN\nhello world\n",
				probe:         curl.Probe{},
			},
			want: false,
		},
		{
			name: "only end token",
			args: args{
				consoleOutput: "hello world\nNV_CURLJSON_END\njj",
				probe:         curl.Probe{},
			},
			want: false,
		},
		{
			name: "token order reversed",
			args: args{
				consoleOutput: "fjsdklNV_CURLJSON_END\nhello world\nNV_CURLJSON_BEGIN\njj",
				probe:         curl.Probe{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokens(tt.args.consoleOutput, tt.args.probe); got != tt.want {
				t.Errorf("get_tokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
