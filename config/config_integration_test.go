// +build integration

package config

import "testing"

func TestGetLicenseID(t *testing.T) {
	c := DefaultConfig()
	if err := c.getLicenseID(); err != nil {
		t.Error("Expected no errors from getLicenseID(), got ", err)
	}

	if c.LicenseID != "upgrade_license" {
		t.Error("Expected license ID to be 'upgrade_license', got ", c.LicenseID)
	}
}
