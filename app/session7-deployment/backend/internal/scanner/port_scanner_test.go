package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

func TestPortScanner_Type(t *testing.T) {
	s := NewPortScanner()
	if s.Type() != model.ScanTypePort {
		t.Fatalf("Type() = %s, want %s", s.Type(), model.ScanTypePort)
	}
}

func TestPortScanner_Scan_UnauthorizedPublicIP(t *testing.T) {
	s := NewPortScanner()
	_, err := s.Scan(&model.Asset{Name: "8.8.8.8", Type: model.TypeIP})
	if err == nil {
		t.Fatal("expected error for unauthorized public ip")
	}
}

func TestPortScanner_Scan_Localhost(t *testing.T) {
	s := NewPortScanner()
	result, err := s.Scan(&model.Asset{Name: "127.0.0.1", Type: model.TypeIP})
	if err != nil {
		t.Fatalf("Scan() unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.TotalScanned <= 0 {
		t.Fatalf("TotalScanned = %d, want > 0", result.TotalScanned)
	}
	if result.ClosedPorts+len(result.OpenPorts) != result.TotalScanned {
		t.Fatalf("open + closed must equal total scanned")
	}
}
