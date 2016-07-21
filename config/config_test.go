package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	d := DefaultConfig()
	if d.SegmentEvent != "health" {
		t.Error("Expected default segment event to be health, got", d.SegmentEvent)
	}

	if d.DCOSClusterIDPath != "/var/lib/dcos/cluster-id" {
		t.Error("Cluster ID path invalid, got", d.DCOSClusterIDPath)
	}
}

func TestFlagParsing(t *testing.T) {
	json := []byte(`{"cluster_uuid": "12345"}`)
	tempAnonJson, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempAnonJson.Name())
	tempAnonJson.Write(json)

	config := []byte(`{"customer_key": "someuser-enterprise-key"}`)
	tempConfig, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempConfig.Name())
	tempConfig.Write(config)

	var (
		verboseConfig, verboseErr = ParseArgsReturnConfig([]string{
			"-v",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		versionConfig, versionErr = ParseArgsReturnConfig([]string{
			"-version",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		testConfig, testConfigErr = ParseArgsReturnConfig([]string{
			"-test",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})
	)

	// -test
	if !testConfig.FlagTest {
		t.Error("Expected test flag to be true, got ", testConfig.FlagTest)
	}
	if testConfigErr != nil {
		t.Error("Expected test error to be nil, got ", testConfigErr)
	}

	// -v
	if verboseConfig.FlagVerbose != true {
		t.Error("Expected true, got ", verboseConfig.FlagVerbose)
	}
	if verboseErr != nil {
		t.Error("Expected nil, got ", verboseErr)
	}

	// -version
	if versionConfig.FlagVersion != true {
		t.Error("Expected true, got ", versionConfig.FlagVersion)
	}
	if versionErr != nil {
		t.Error("Expected nil, got ", versionErr)
	}
}

func TestGetClusterID(t *testing.T) {
	file := []byte(`12345`)
	tempAnon, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempAnon.Name())

	tempAnon.Write(file)
	c := DefaultConfig()
	c.DCOSClusterIDPath = tempAnon.Name()

	if err := c.getClusterID(); err != nil {
		t.Error("Expected no errors from getClusterID(), got ", err)
	}

	if c.ClusterID != "12345" {
		t.Error("Expected cluster ID to be 12345, got ", c.ClusterID)
	}
}

func TestGetExternalConfig(t *testing.T) {

	noEntConfig := []byte(`
{
"customer_key": "",
"gen_provider": "onprem"	
}`)
	noEntFile, _ := ioutil.TempFile(os.TempDir(), "")

	defer os.Remove(noEntFile.Name())

	noEntFile.Write(noEntConfig)
	noEntC := DefaultConfig()
	noEntC.SignalServiceConfigPath = noEntFile.Name()

	if err := noEntC.getExternalConfig(); err != nil {
		t.Error("Expected no errors with no enterprise in config, got ", err)
	}

	if noEntC.CustomerKey != "" {
		t.Error("Expected customer ID to be empty, got ", noEntC.CustomerKey)
	}

	if noEntC.GenProvider != "onprem" {
		t.Error("Expected onprem, got ", noEntC.GenProvider)
	}
}

func TestTryLoadingCert(t *testing.T) {
	mockCertData := []byte(`
-----BEGIN CERTIFICATE-----
MIIFtTCCA52gAwIBAgIJAMDDrJnHaFjRMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTYwNjA0MjExMTI4WhcNMzYwNTMwMjExMTI4WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAwRMKKNwqPpWkvglFJ9X1G8KSz2yOitlbzQG/0oiYUPhBkI7OWud9nJ/j
cJiV/GycmBQm+RmZdFXSdFcyDupBrMKoddGtYZR1NWCwsKesw5OczOvI9WHemLGa
meBNH3plzaXdFp/0H/EwLp7DWypIaw9KoC5i6nRO42EvvwZ/Vk77LrFdOEFEv4q7
LcJMsO/AdSfLXshj0xC7NYCo2ZPS5q+OQfvTcgbyzYMYg0u7yBvX3svwZlY92OYX
FtEdFXYmLFovk7cgSgpizaDcqYk/Yn/iYLKsKynNxwCFY6FjzVB+yGquEUtr+Om+
FhXqPlts4dXforvoUddn2sRHv6ctYmWIbMvS6vxNyISuHjOQ7WqsOHOZOSUSEsW4
h/MXKd7UUkRjjPrvDl5ZMJnkbQndourB+C6v6snSqwmCHop+ZUizsn5fhnbn6xh8
QV1ElMS8p51lEIHxbN+evxpNTPfXVdO9ZXasJZDcfWr9RfTFct18ep3kqCt7ccWC
FR+QLZ+4ga9BEzOWkNxbZa1nouWWbbWj/P6tmVlCEIWIQGu5YfxvIY8m6zrANaQ/
f8Lr7XHAV92KaSUWdKVAwOBV3NExWoZ2EHUloYmwrVwKWa7qDCAiDXniy9QyuDfh
eNbQoUDLvYn7AZo82bzN99ymgi9FiJMck66r8oamyuplARP8F7MCAwEAAaOBpzCB
pDAdBgNVHQ4EFgQUhHAUX2lq3kBZdcg3wqEbCxY2/hYwdQYDVR0jBG4wbIAUhHAU
X2lq3kBZdcg3wqEbCxY2/hahSaRHMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpT
b21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGSCCQDA
w6yZx2hY0TAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQCjC/+PsxMi
b2Fc8W78u5dXRlzCZqMUuHSlhRD8Q3+nbV3aslRUOsAwo9U7PHJepD2BgzaQwAbW
H8nwnOERzCPH6jIFRJR6IPM8in4usc5Q0AEJxprxF6vNz4/27bDXSzvUgBroawHK
JDflt44a54kVBu9LHVANiW4Ydn2fFsb44um2BdgXLUYyBL6zIu2iuEq7g41QKUP9
1AFGXDXflYLvy52SdP7krtTs5G66XfOK2lAXmDla1aJ5yCmkomI2WaOziwWkg4A8
NLBeRAY9v8ppOgnHRCxFs6y4zAtYusY6MU3txR4P8aqFLqyEllRbyeHygOp3Npun
MXuCZEho8KDSgteIqCWXOaRqJZW4CgDDb2/dX/NFO3MLGZff4Nno25w/qpbqIFNT
ZOv25WrKzBJnmEd8KFkNWNJOr+BMH3EvVRZaSM2jHTYvTKNsF9Rm5ubY/JSRpMzi
qrNebMwmr+uDKx2PFQ/1hQGXB/B7ognF4f1kogPa59qHu/tpo4awmu9GvO3H/BkC
2kvyQ2N+o1UN2tAkBwZlmNPZWzn7wtP7Gf08tPDuFBfJAIJqxfQwZOA61juzZhsp
6/ft+Mu6u35mX2stebsvSiipsli/AuDKaMHgZW7SJGX/JoKkjQwdxpupUSwBzAUm
NfOBa2wja44Izn/W58mcHkIvK61BCWGKTg==
-----END CERTIFICATE-----`)
	mockCert, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(mockCert.Name())
	mockCert.Write(mockCertData)

	mockC := Config{
		CACertPath: mockCert.Name(),
	}

	if err := mockC.tryLoadingCert(); err != nil {
		t.Error("Expected no errors loading cert, got", err)
	}

	if mockC.CAPool == nil {
		t.Error("Expected CAPool, got", mockC.CAPool)
	}
}
