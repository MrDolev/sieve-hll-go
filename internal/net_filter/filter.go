package netfilter

import (
	"fmt"
	"net/netip"
	"strings"
)

type NetFilterI interface {
	Filter(ip string) bool
}

type NetFilter struct {
	cidr []string
}

func NewNetFilter(cidr []string) *NetFilter {
	return &NetFilter{
		cidr: cidr,
	}
}

// Filter reports whether ip falls within at least one of the configured CIDR
// ranges, meaning it must be excluded from collection. It returns false when
// ip is not a valid address or matches none of them.
func (nt *NetFilter) Filter(ip string) bool {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	for _, mask := range nt.cidr {
		prefix, err := netip.ParsePrefix(mask)
		if err != nil {
			continue
		}
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

// NormalizeCIDR validates every entry of cidr and returns them in canonical
// form. A bare address is turned into a single-host prefix (/32 or /128).
// Empty entries are skipped.
func NormalizeCIDR(cidr []string) ([]string, error) {
	results := make([]string, 0, len(cidr))
	for _, entry := range cidr {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if addr, err := netip.ParseAddr(entry); err == nil {
			addr = addr.Unmap()
			results = append(results, netip.PrefixFrom(addr, addr.BitLen()).String())
			continue
		}
		prefix, err := netip.ParsePrefix(entry)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR %q: %w", entry, err)
		}
		results = append(results, prefix.Masked().String())
	}
	return results, nil
}
