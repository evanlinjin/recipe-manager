package config

import (
	"testing"
	"os"
)

func TestGetNetworkConfig(t *testing.T) {
	c, e := GetNetworkConfig()
	if e != nil {
		t.Fatal(e)
	}
	if c.Domain == "" {
		t.Fatal("Empty Domain")
	}
	if c.Port == "" {
		t.Fatal("Empty Port")
	}
	if c.Mongo == "" {
		t.Fatal("Empty Mongo")
	}
	if c.SSLCertPath == "" {
		t.Fatal("Empty SSLCertPath")
	}
	if c.SSLKeyPath == "" {
		t.Fatal("Empty SSLKeyPath")
	}
	os.Remove(TempSSLCertPath)
	os.Remove(TempSSLKeyPath)
}
