package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

func TestWHOISScanner_Type(t *testing.T) {
	s := NewWHOISScanner()
	if s.Type() != model.ScanTypeWHOIS {
		t.Fatalf("Type() = %s, want %s", s.Type(), model.ScanTypeWHOIS)
	}
}

func TestWHOISScanner_Scan_InvalidAssetType(t *testing.T) {
	s := NewWHOISScanner()
	_, err := s.Scan(&model.Asset{Name: "127.0.0.1", Type: model.TypeIP})
	if err == nil {
		t.Fatal("expected error for non-domain asset")
	}
}

func TestWHOISScanner_GetWHOISServer(t *testing.T) {
	s := NewWHOISScanner()
	server := s.getWHOISServer("example.com")
	if server == "" {
		t.Fatal("expected non-empty whois server")
	}
}
