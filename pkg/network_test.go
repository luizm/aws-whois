package whois

import "testing"

func TestIsIP(t *testing.T) {
	validaIP := "10.1.1.1"
	invalidIP := "10.1.1.256"
	if !IsIP(validaIP) {
		t.Errorf("the valid ip was considered invalid %s", validaIP)
	}

	if IsIP(invalidIP) {
		t.Errorf("the invalid ip was considered valid %s", invalidIP)
	}
}

func TestPrivateIP(t *testing.T) {
	privateIP, _ := isPrivateIP("10.1.1.1")
	publicIP, _ := isPrivateIP("8.8.8.8")

	if !privateIP {
		t.Error("the private ip was considered public")
	}

	if publicIP {
		t.Error("the private ip was considered private")
	}
}

func TestResolvDNS(t *testing.T) {
	ips, err := ResolvDNS("apple.com")
	if err != nil {
		t.Error("failed to resolve DNS name")
	}

	if !IsIP(ips[0]) {
		t.Error("invalid result in the DNS query")
	}
}
