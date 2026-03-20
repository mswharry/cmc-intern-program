package scanner

import (
	"context"
	"testing"
	"time"

	"mini-asm/internal/model"
)

func TestSubdomainScanner_Type(t *testing.T) {
	s, err := NewSubdomainScanner()
	if err != nil {
		t.Fatalf("NewSubdomainScanner() error: %v", err)
	}
	if s.Type() != model.ScanTypeSubdomain {
		t.Fatalf("Type() = %s, want %s", s.Type(), model.ScanTypeSubdomain)
	}
}

func TestSubdomainScanner_Scan_InvalidAssetType(t *testing.T) {
	s, err := NewSubdomainScanner()
	if err != nil {
		t.Fatalf("NewSubdomainScanner() error: %v", err)
	}
	_, scanErr := s.Scan(&model.Asset{Name: "127.0.0.1", Type: model.TypeIP}, context.Background())
	if scanErr == nil {
		t.Fatal("expected error for non-domain asset")
	}
}

func TestSubdomainScanner_Scan_LocalhostDomain(t *testing.T) {
	s, err := NewSubdomainScanner()
	if err != nil {
		t.Fatalf("NewSubdomainScanner() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results, scanErr := s.Scan(&model.Asset{Name: "localhost", Type: model.TypeDomain}, ctx)
	if scanErr != nil {
		t.Fatalf("Scan() unexpected error: %v", scanErr)
	}
	if results == nil {
		t.Fatal("expected non-nil subdomain results")
	}
}
