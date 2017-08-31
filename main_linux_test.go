package main

import (
	"bytes"
	"fmt"
	"net/mail"
	"os"
	"os/exec"
	"testing"
)

var testMail = []byte(`Return-Path: <example@example.com>
Received: from mail.example.com (localhost [127.0.0.1])
	by mail.example.com (Postfix) with ESMTP id 2CF6CE0EF4
	for <user@example.com>; Thu, 31 Aug 2017 01:06:32 +0200 (CEST)
Date: Thu, 31 Aug 2017 01:06:30 +0200
From: example@example.com
To: user@example.com
Message-ID: <59a744f6.dpQIYZRbHlm4eKow%example@example.com>
User-Agent: Heirloom mailx 12.5 7/5/10
MIME-Version: 1.0
Message-ID-Hash: CAM2MJNGYWXL7UWVJLIV47ZH2SENJGWM
X-Message-ID-Hash: CAM2MJNGYWXL7UWVJLIV47ZH2SENJGWM
Precedence: list
Subject: =?iso-8859-1?q?=5BTest=5D_test=FC-SMIME?=
Content-Type: text/plain; charset="iso-8859-1"
Content-Transfer-Encoding: quoted-printable

This is the body of a mail.
`)

var cert = []byte(`-----BEGIN TRUSTED CERTIFICATE-----
MIIE/jCCAuYCAQEwDQYJKoZIhvcNAQELBQAwRTELMAkGA1UEBhMCQVUxEzARBgNV
BAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0
ZDAeFw0xNzA4MzEyMDA3NTJaFw0yNzA4MjkyMDA3NTJaMEUxCzAJBgNVBAYTAkFV
MRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRz
IFB0eSBMdGQwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDENdoac97E
qFxqbuBOJXIpc5C+2R2LB95Z/xN4+FxTBs5D9rXvDZVGVe9FGDFcJTuuMMhT2GPx
QUwfM2LmvEPy1f+2b+tIpIasESUcf+LtT459t0ogra8J2nX7Tti71z0QuZZBvBzW
3dcYa4B0J9ABIMVGnJ43Nw0orW6vgvYmLv/m9ai09ZXSD8/210Uc+GSESxWTlaQY
psi/UiFAGwFCEEyaxCafq6c5UkzsSXI5jta+Fg/0ISmW7DSXL7m2HCXQaCTR7CaN
DVwrCniyal3wMj0G+m5cHJEC9/hqEPeLNOaXZXbBbeF1LQvQLiuGI6mmKmQ3Q3q1
AVsckJyAJUB50026tell/T+wmBtIkFR/5zXt+Al0xDweIpwLjAb1fNGm2flkwgC5
AZUr6pGJBdHe7bzU/BnKVPSgivgSt3fBWE5g4BW2vCK+TvxYmAeOWYBSx4mmKK5j
0RZ6P+UvBP62XxpAtdrGhM1XNZY6CoaVjvPZkZ6GypPUl4g2EPgkLlfl5Gu79KWp
1agWSkc47SqE45UR8CVdeqWp8ABMfAmB+GMiWZzs9oM23sLwcbETNOngfNyCi0WB
z3RP4P2rIpWuTz1bXguyrf8dHwB8IiyDQttzf8yG1QG9a99bSR1a4irR0TGcd7M2
mwo8FA2HOwPBXPSMtADW7yZhlvPMAmNsswIDAQABMA0GCSqGSIb3DQEBCwUAA4IC
AQBRK1EPxkyWv8tnUM7XFZsBqOJpdZtzYFXqrekpr6JjYcfia+dHl72//cJNz7xE
UhdA0tUjQuno8Z34MeERwGfrO2XQUWM5KvrjDB0gryBxXmg4ZaUAKBT/u/i81ebL
bXkE0Q83VwahbX+t4hZl0hi7zDXQn8teFELYWMf4Yi9E2uYb8nt6MlKI2Db+NNz9
EJ3EVyxk6t58fXhBGYJJvRF0mDxNUz4R5mgFiN8xZRdIKs4V7YTay0tSJlTfMGPa
PmQiNFu6j+FP1hNNbSrekDR7w174OjQzjkULWygMPLolDn1yW5DoSpDFNcCxe3jb
MtGBFxgBMNwfnSG5LbLIu0qSUenupQ8VFKGhmkf2bmKYGh+9CWyCrpXtuT4vT4AY
mXHU16/Jtp6/dCDgMESCdhi9VxOLOJ0/i87uqtFzR8lllSojYEEw8VrAeWsu4rYg
fazvBbTQmsY2pM+xthm3qHTg4mLL1pfg5cjJC+V5ZAIxpYQ0ydlGz8RIdkqCEh0Y
lamq+COEcgxnBR87HRlqK50Yuffh1mE+AOb15cNc+K3WuwJ8oP6Op4E8GBiYAWQU
NYQy2xHJS8SiPfRvGjHhfcdQ6zg8U5JU4Rx+kWbChhQSUCUbfWT6XTLo8mD4TxQ1
D55EpsWYKZtByXKGz/Z646XvfwiQVUiFHL8iEglND/jvZzA1MAoGCCsGAQUFBwME
oBQGCCsGAQUFBwMCBggrBgEFBQcDAQwRU2VsZiBTaWduZWQgU01JTUU=
-----END TRUSTED CERTIFICATE-----
`)

var key = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)

func writeFile(path string, b []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create file: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("couldn't write file:  %v", err)
	}
	return nil
}

func TestEncryption(t *testing.T) {
	var certdir = os.TempDir()
	rcpt := "user@example.com"
	cfile := certdir + "/" + rcpt + ".crt"
	kfile := certdir + "/" + rcpt + ".key"
	if err := writeFile(cfile, cert); err != nil {
		t.Fatalf("couldn't create file %v: %v", cfile, err)
	}
	if err := writeFile(kfile, key); err != nil {
		t.Fatalf("couldn't create file %v: %v", kfile, err)
	}
	defer os.Remove(cfile)
	defer os.Remove(kfile)

	var buf bytes.Buffer
	_, err := buf.Write(testMail)
	if err != nil {
		t.Fatalf("couldn't write testMail to buffer: %v", err)
	}

	emsg, err := encryptMail(certdir, "example@example.org", rcpt, "test-SMIME", buf)
	if err != nil {
		t.Fatalf("failed to encrypt mail: %v, %s", err, emsg)
	}

	c := exec.Command("openssl", "smime", "-decrypt", "-recip", cfile, "-inkey", kfile)
	// Pipe msg to openssl.
	stdin, err := c.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	// Write to stdin.
	go func() {
		defer stdin.Close()
		stdin.Write(emsg)
	}()

	cmsg, err := c.Output()
	if err != nil {
		t.Fatal("couldn't decrypt msg: ", err)
	}

	m1, err := mail.ReadMessage(bytes.NewReader(testMail))
	m2, err := mail.ReadMessage(bytes.NewReader(cmsg))

	var b1, b2 []byte
	m1.Body.Read(b1)
	m2.Body.Read(b2)

	if same := bytes.Compare(b1, b2); same != 0 {
		t.Fatalf("encrypted mail should be \n%q \n\n got\n%q", b1, b2)
	}
}
