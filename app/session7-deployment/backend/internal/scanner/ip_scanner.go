package scanner

import (
	"fmt"
	"net"
	"strings"
	"time"

	"mini-asm/internal/model"
)

type IPScanner struct {
	timeout time.Duration
}

func NewIPScanner() *IPScanner {
	return &IPScanner{timeout: 5 * time.Second}
}

func (s *IPScanner) Type() model.ScanType {
	return model.ScanTypeIP
}

func (s *IPScanner) Scan(asset *model.Asset) (*model.IPScanResult, error) {
	if asset.Type != model.TypeIP {
		return nil, fmt.Errorf("ip scan requires ip asset, got: %s", asset.Type)
	}
	ip := net.ParseIP(asset.Name)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip address: %s", asset.Name)
	}

	rdns := ""
	names, _ := net.LookupAddr(asset.Name)
	if len(names) > 0 {
		rdns = strings.TrimSuffix(names[0], ".")
	}

	return &model.IPScanResult{
		IPAddress: asset.Name,
		Geolocation: model.GeoLocation{
			Country:     "Local",
			CountryCode: "LO",
			City:        "Localhost",
			Region:      "Loopback",
			Latitude:    0,
			Longitude:   0,
			ISP:         "Local Network",
			Org:         "Training",
		},
		ASN: model.ASNInfo{
			Number:      0,
			Name:        "LOCAL",
			Description: "Private/local network",
		},
		ReverseDNS: rdns,
	}, nil
}
