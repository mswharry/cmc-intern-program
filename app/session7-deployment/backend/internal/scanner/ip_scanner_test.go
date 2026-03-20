package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

func TestIPScanner_Type(t *testing.T) {
	s := NewIPScanner()
	if s.Type() != model.ScanTypeIP {
		t.Fatalf("Type() = %s, want %s", s.Type(), model.ScanTypeIP)
	}
}

func TestIPScanner_Scan_InvalidAssetType(t *testing.T) {
	s := NewIPScanner()
	_, err := s.Scan(&model.Asset{Name: "example.com", Type: model.TypeDomain})
	if err == nil {
		t.Fatal("expected error for non-ip asset")
	}
}

func TestIPScanner_Scan_ValidLocalIP(t *testing.T) {
	s := NewIPScanner()
	result, err := s.Scan(&model.Asset{Name: "127.0.0.1", Type: model.TypeIP})
	if err != nil {
		t.Fatalf("Scan() unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.IPAddress != "127.0.0.1" {
		t.Fatalf("IPAddress = %s, want 127.0.0.1", result.IPAddress)
	}
}
