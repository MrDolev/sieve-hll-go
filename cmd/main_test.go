package main

import (
	"reflect"
	"strings"
	"testing"

	netfilter "github.com/mrdolev/sieve-hll-go/internal/net_filter"
)

func TestCIDRFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		want     []string
		wantErr  bool
		filtered []string
		allowed  []string
	}{
		{
			name:    "unset filters nothing",
			env:     "",
			want:    []string{},
			allowed: []string{"10.0.0.1", "192.168.0.1", "8.8.8.8"},
		},
		{
			name:     "private classes",
			env:      "10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16",
			want:     []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
			filtered: []string{"10.0.0.0", "10.255.255.255", "172.16.0.1", "172.31.255.255", "192.168.1.50"},
			allowed:  []string{"11.0.0.1", "172.15.255.255", "172.32.0.1", "192.169.0.1", "8.8.8.8"},
		},
		{
			name:     "single hosts and loopback",
			env:      "8.8.8.8,127.0.0.0/8",
			want:     []string{"8.8.8.8/32", "127.0.0.0/8"},
			filtered: []string{"8.8.8.8", "127.0.0.1", "127.255.255.254"},
			allowed:  []string{"8.8.4.4", "128.0.0.1"},
		},
		{
			name:    "wildcard is rejected",
			env:     "10.0.*.*",
			wantErr: true,
		},
		{
			name:    "one invalid entry rejects the whole list",
			env:     "10.0.0.0/8,192.168.0.0/33",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CIDR", tt.env)

			cidr, err := getEnvCIDR("CIDR")
			if (err != nil) != tt.wantErr {
				t.Fatalf("getEnvCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(cidr, tt.want) {
				t.Fatalf("CIDR = %v, want %v", cidr, tt.want)
			}

			filter := netfilter.NewNetFilter(cidr)
			for _, ip := range tt.filtered {
				if !filter.Filter(ip) {
					t.Errorf("%s should be filtered by %v", ip, cidr)
				}
			}
			for _, ip := range tt.allowed {
				if filter.Filter(ip) {
					t.Errorf("%s should not be filtered by %v", ip, cidr)
				}
			}
		})
	}
}

func TestTrustedProxiesFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		want    []string
		wantErr bool
	}{
		{
			name: "unset trusts no proxy",
			env:  "",
			want: []string{},
		},
		{
			name: "local development hosts",
			env:  "127.0.0.1, ::1, 192.168.65.1",
			want: []string{"127.0.0.1/32", "::1/128", "192.168.65.1/32"},
		},
		{
			name: "proxy subnets",
			env:  "10.0.0.0/8,172.16.0.0/12",
			want: []string{"10.0.0.0/8", "172.16.0.0/12"},
		},
		{
			name: "ipv6 proxy subnet",
			env:  "2001:db8::/32",
			want: []string{"2001:db8::/32"},
		},
		{
			name: "only spaces trusts no proxy",
			env:  " , ",
			want: []string{},
		},
		{
			name:    "invalid entry is rejected",
			env:     "10.0.0.0/8,proxy.local",
			wantErr: true,
		},
		{
			name:    "ipv4 catch-all is rejected",
			env:     "0.0.0.0/0",
			wantErr: true,
		},
		{
			name:    "ipv6 catch-all is rejected",
			env:     "10.0.0.0/8,::/0",
			wantErr: true,
		},
		{
			name:    "catch-all with host bits is rejected",
			env:     "10.0.0.1/0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TRUSTED_PROXIES", tt.env)

			got, err := getTrustedProxies()
			if (err != nil) != tt.wantErr {
				t.Fatalf("getTrustedProxies() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), "TRUSTED_PROXIES") {
					t.Errorf("error %q should name the variable", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TRUSTED_PROXIES = %v, want %v", got, tt.want)
			}
		})
	}
}
