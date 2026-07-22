package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	if cfg.Port != "3000" {
		t.Fatalf("expected 3000, got: %s", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected info, got: %s", cfg.LogLevel)
	}
	if cfg.TranscodeWorkers != 2 {
		t.Fatalf("expected 2, got: %d", cfg.TranscodeWorkers)
	}
}

func TestLoad_EnvVarOverridesDefault(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Fatalf("expected 8080, got: %s", cfg.Port)
	}
}

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		vals []string
		want string
	}{
		{[]string{"a", "b", "c"}, "a"},
		{[]string{"", "b", "c"}, "b"},
		{[]string{"", "", "c"}, "c"},
		{[]string{"", "", ""}, ""},
		{[]string{}, ""},
	}
	for _, tt := range tests {
		got := firstNonEmpty(tt.vals...)
		if got != tt.want {
			t.Fatalf("firstNonEmpty(%v) = %q; want %q", tt.vals, got, tt.want)
		}
	}
}

func TestGetDurationEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_DURATION", "1h")
	defer os.Unsetenv("TEST_DURATION")

	d := getDurationEnv("TEST_DURATION", 0)
	if d.Hours() != 1 {
		t.Fatalf("expected 1h, got: %v", d)
	}
}

func TestGetInt64Env(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	n := getInt64Env("TEST_INT", 0)
	if n != 42 {
		t.Fatalf("expected 42, got: %d", n)
	}
}

func TestGetSliceEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_SLICE", "http://localhost:3000")
	defer os.Unsetenv("TEST_SLICE")

	s := getSliceEnv("TEST_SLICE", []string{"*"})
	if len(s) != 1 || s[0] != "http://localhost:3000" {
		t.Fatalf("expected [http://localhost:3000], got: %v", s)
	}
}

func TestLoad_JWTSecretDefault(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	if cfg.JWTSecret != "change-me-in-production" {
		t.Fatalf("expected change-me-in-production, got: %s", cfg.JWTSecret)
	}
}
