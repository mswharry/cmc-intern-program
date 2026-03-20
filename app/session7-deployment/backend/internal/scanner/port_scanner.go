package scanner

import (
	"fmt"
	"log"
	"net"
	"time"

	"mini-asm/internal/model"
)

// PortScanner performs port scanning on IP addresses or hostnames.
//
// SCAN CATEGORY: ACTIVE - requires explicit authorization.
// For this training project, the scanner only allows localhost/private addresses.
type PortScanner struct {
	timeout       time.Duration
	maxWorkers    int
	commonPorts   []int
	authorizedIPs map[string]bool
}

// NewPortScanner creates a new port scanner with safety restrictions.
func NewPortScanner() *PortScanner {
	log.Println("active port scan initialized with safety checks")

	return &PortScanner{
		timeout:    5 * time.Second,
		maxWorkers: 100,
		commonPorts: []int{
			21,
			22,
			23,
			25,
			53,
			80,
			110,
			143,
			443,
			445,
			3306,
			3389,
			5432,
			5900,
			8080,
			8443,
		},
		authorizedIPs: map[string]bool{
			"127.0.0.1": true,
			"localhost": true,
			"::1":       true,
		},
	}
}

// Type returns the scan type identifier.
func (s *PortScanner) Type() model.ScanType {
	return model.ScanTypePort
}

// Scan performs port scanning on a target IP address/domain.
func (s *PortScanner) Scan(asset *model.Asset) (*model.PortScanResult, error) {
	if asset.Type != model.TypeIP && asset.Type != model.TypeDomain {
		return nil, fmt.Errorf("port scan requires ip or domain asset, got: %s", asset.Type)
	}

	target := asset.Name
	if !s.isAuthorized(target) {
		return nil, fmt.Errorf("unauthorized port scan target: %s (only localhost/private IPs are allowed)", target)
	}

	log.Printf("active scan started: scan_type=port target=%s", target)

	openPorts := make([]model.OpenPort, 0)
	closedPorts := 0
	start := time.Now()

	for _, port := range s.commonPorts {
		if s.scanPort(target, port) {
			openPorts = append(openPorts, model.OpenPort{
				Port:     port,
				Protocol: "tcp",
				State:    "open",
				Service:  s.detectService(port),
				Version:  "",
			})
		} else {
			closedPorts++
		}
	}

	result := &model.PortScanResult{
		IPAddress:      target,
		OpenPorts:      openPorts,
		ClosedPorts:    closedPorts,
		TotalScanned:   len(s.commonPorts),
		ScanDurationMS: int(time.Since(start).Milliseconds()),
	}

	log.Printf("active scan completed: scan_type=port target=%s open=%d closed=%d", target, len(openPorts), closedPorts)

	_ = s.maxWorkers
	return result, nil
}

func (s *PortScanner) isAuthorized(target string) bool {
	if s.authorizedIPs[target] {
		return true
	}

	ips, err := net.LookupIP(target)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if s.authorizedIPs[ip.String()] || ip.IsPrivate() || ip.IsLoopback() {
			return true
		}
	}

	return false
}

func (s *PortScanner) scanPort(target string, port int) bool {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func (s *PortScanner) detectService(port int) string {
	switch port {
	case 21:
		return "ftp"
	case 22:
		return "ssh"
	case 25:
		return "smtp"
	case 53:
		return "dns"
	case 80:
		return "http"
	case 110:
		return "pop3"
	case 143:
		return "imap"
	case 443:
		return "https"
	case 445:
		return "smb"
	case 3306:
		return "mysql"
	case 3389:
		return "rdp"
	case 5432:
		return "postgresql"
	case 5900:
		return "vnc"
	case 8080:
		return "http-alt"
	case 8443:
		return "https-alt"
	default:
		return "unknown"
	}
}
