package netfilter

import (
	"reflect"
	"testing"
)

func TestNormalizeCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidr    []string
		want    []string
		wantErr bool
	}{
		{
			name: "keeps valid cidr and trims spaces",
			cidr: []string{" 10.0.0.0/8", "192.168.0.0/16 "},
			want: []string{"10.0.0.0/8", "192.168.0.0/16"},
		},
		{
			name: "bare ipv4 becomes /32",
			cidr: []string{"10.0.0.1"},
			want: []string{"10.0.0.1/32"},
		},
		{
			name: "bare ipv6 becomes /128",
			cidr: []string{"2001:db8::1"},
			want: []string{"2001:db8::1/128"},
		},
		{
			name: "masks host bits",
			cidr: []string{"192.168.1.50/24"},
			want: []string{"192.168.1.0/24"},
		},
		{
			name: "skips empty entries",
			cidr: []string{"10.0.0.0/8", "", " "},
			want: []string{"10.0.0.0/8"},
		},
		{
			name: "empty list",
			cidr: []string{},
			want: []string{},
		},
		{
			name:    "rejects wildcard",
			cidr:    []string{"10.0.*.*"},
			wantErr: true,
		},
		{
			name:    "rejects invalid entry",
			cidr:    []string{"10.0.0.0/8", "not-a-cidr"},
			wantErr: true,
		},
		{
			name:    "rejects out of range prefix",
			cidr:    []string{"10.0.0.0/33"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeCIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetFilter_Filter(t *testing.T) {
	private := []string{"172.16.0.0/12", "192.168.0.0/16"}
	tests := []struct {
		name string
		cidr []string
		ip   string
		want bool
	}{
		{
			name: "ip in second cidr",
			cidr: private,
			ip:   "192.168.0.1",
			want: true,
		},
		{
			name: "ip in first cidr",
			cidr: private,
			ip:   "172.31.255.255",
			want: true,
		},
		{
			name: "ip outside every cidr",
			cidr: private,
			ip:   "10.0.0.1",
			want: false,
		},
		{
			name: "no cidr configured",
			cidr: []string{},
			ip:   "10.0.1.1",
			want: false,
		},
		{
			name: "exact /32 match",
			cidr: []string{"10.0.0.1/32"},
			ip:   "10.0.0.1",
			want: true,
		},
		{
			name: "/32 does not match neighbour",
			cidr: []string{"10.0.0.1/32"},
			ip:   "10.0.0.2",
			want: false,
		},
		{
			name: "invalid entry does not skip the following ones",
			cidr: []string{"not-a-cidr", "192.168.0.0/16"},
			ip:   "192.168.0.1",
			want: true,
		},
		{
			name: "invalid ip is not filtered and does not panic",
			cidr: private,
			ip:   "not-an-ip",
			want: false,
		},
		{
			name: "ipv4-mapped ipv6 matches ipv4 cidr",
			cidr: private,
			ip:   "::ffff:192.168.0.1",
			want: true,
		},
		{
			name: "ipv6 cidr",
			cidr: []string{"2001:db8::/32"},
			ip:   "2001:db8::1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			nt := NewNetFilter(tt.cidr)
			if got := nt.Filter(tt.ip); got != tt.want {
				t.Errorf("NetFilter.Filter(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
