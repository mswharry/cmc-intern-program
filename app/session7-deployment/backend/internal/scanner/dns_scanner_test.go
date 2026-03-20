package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

func TestDNSScanner_Type(t *testing.T) {
	s := NewDNSScanner()
	if s.Type() != model.ScanTypeDNS {
		t.Fatalf("Type() = %s, want %s", s.Type(), model.ScanTypeDNS)
	}
}

func TestDNSScanner_Scan_InvalidAssetType(t *testing.T) {
	s := NewDNSScanner()
	_, err := s.Scan(&model.Asset{Name: "127.0.0.1", Type: model.TypeIP})
	if err == nil {
		t.Fatal("expected error for non-domain asset")
	}
}

func TestDNSScanner_Scan_LocalhostDomain(t *testing.T) {
	s := NewDNSScanner()
	records, err := s.Scan(&model.Asset{Name: "localhost", Type: model.TypeDomain})
	if err != nil {
		t.Fatalf("Scan() unexpected error: %v", err)
	}
	if records == nil {
		t.Fatal("expected non-nil records slice")
	}
}
