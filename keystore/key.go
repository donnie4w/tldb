// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb

package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/donnie4w/tldb/log"
	"github.com/donnie4w/tldb/util"
)

var defaultPriKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCVXHX5DUrg2cNQTwnAIrJC4r1u3MewXBIaUA6+yrWO1EAh4ZwC
8n86WzM2rqOmU/qAmTScAeAfcRJkq0utT+ETW97W3sHAbCCpvunDRP0ZrEhlqzZ2
dEtdwB8kPavQLMD0usXAyIRJZOYdzYOXvWJGzU3T4SAbTHdThrTII/nA3QIDAQAB
AoGAFgGhgChtN+Pd2x9KGH0ENsahkowE//8Qy7+v7HyBc6HiMRvEmMqR5E87pHrm
scL9zaTFE5dTJk7Knvp+E/MI/k+S2Z02tZajohAR07x4rZSIGZ0s8CsVxwN2J4UO
zqIrCULNIrQxQfUDodFA3v28c5UeKMNb7tRW8fsk14+TXQECQQDFvLu9TYa/6CLJ
N0Ap2SVjO9TqzJX/Vpi8QxQxFUNaWIapees5yA8T85l8dcim+p56NZNCfKRK5BpG
7sZgecAZAkEAwV675divPojkm6paASXPB8j0uJq6tADNrZyRIWr49w5pCNraZotI
XECZZY5VK0+n29m9gQhdVdYUE7BNJxuPZQJAeacj2dNYk7i9rg3P6+8skWC+HbbA
kdc1IJ4kTg5G4c6VCq93iJUMsbmtNGVCjXijB4zujHkimvC7OeitI63EAQJAArkt
1kfd9/h/l72ndNqudsKax7rOJFjajLZmyNyz0u7uBcTnTIhrpXj3cBm4E1sU1yDS
7W1Luzi/oaNbAtD9jQJAJVK9BGnDQLZfiJGmceTatCayrKMVG3xwQU/KPSHAWsXp
KM1gVd5HBoqKdEQBHRhSpbvEMyh6tl1oLklu2haavA==
-----END RSA PRIVATE KEY-----`

var defaultPubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCVXHX5DUrg2cNQTwnAIrJC4r1u
3MewXBIaUA6+yrWO1EAh4ZwC8n86WzM2rqOmU/qAmTScAeAfcRJkq0utT+ETW97W
3sHAbCCpvunDRP0ZrEhlqzZ2dEtdwB8kPavQLMD0usXAyIRJZOYdzYOXvWJGzU3T
4SAbTHdThrTII/nA3QIDAQAB
-----END PUBLIC KEY-----`

var ServerCrt = `-----BEGIN CERTIFICATE-----
MIIDUTCCAjkCFHzjAe8MV9uUQHj5r3giJmSeazYDMA0GCSqGSIb3DQEBCwUAMGUx
CzAJBgNVBAYTAkNOMQswCQYDVQQIDAJHRDELMAkGA1UEBwwCR1oxDTALBgNVBAoM
BHRsZGIxCjAIBgNVBAMMAQkxITAfBgkqhkiG9w0BCQEWEmRvbm5pZTR3QGdtYWls
LmNvbTAeFw0yMzA1MjEwMjU2NTRaFw0zMzA1MTgwMjU2NTRaMGUxCzAJBgNVBAYT
AkNOMQswCQYDVQQIDAJHRDELMAkGA1UEBwwCR1oxDTALBgNVBAoMBHRsZGIxCjAI
BgNVBAMMAQkxITAfBgkqhkiG9w0BCQEWEmRvbm5pZTR3QGdtYWlsLmNvbTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALwdkMl+T6+a/N7KucL/KPK4mV2V
1QiOzsBgef4ih95AIRiHypwb0lZxJelsfm2alrtt2jMVINvi4Od92hY/WnhFwtd6
TK5+STIBfgtiPZpImKN0AGPegF3sRBURwcBPyVdcOxYmxjrUuAb0Kg/IBfFyTTe2
Op9z6bGkWfviND9BXFD79OQ4xkdkNZsdR4dQyW50YNOQLCGwVxxDseE40svW/zzu
01QOaUIIExzXnNPUHaVwbJpAq0dPALyooEhZWO3ximMLtit6sMJnlm5aluf3lqWk
lNIsFUnQEmRtXewr+4FwXaaOH/nemg9bgy2jfXBNGuosjt5H/mSHR83eCnMCAwEA
ATANBgkqhkiG9w0BAQsFAAOCAQEAFD1jIM5BHt0ryt9QdGUrRs966b0XeUKOcq+m
CKabLWBj4ucnUQWtfwjJ0GovRreCuIxXkK4uFYO9ov6m4IgINuzWx2SFShor0xOR
aK/ZWf+yZ8465fhHetOYRVBLZYiziZBDC7hYiVxGWqEHllP6KxG4AWswNUr8t0Xv
mN6MTVFQ1WzGGJZSBrHu9hRgg7PAKt5puEwBUWIdH7nwMBI+11B+aUXQrViGKt3/
ucqMkEcTTwn9voJsBMmcheglnQTtEupLHAQU838mZinodO/12Gzt4cyNvmcZnrSR
CHMyOh+50rl4ED84Me6dbzcEmwDHWgH1Afc79onPJzie7hPD5g==
-----END CERTIFICATE-----`

var ServerKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvB2QyX5Pr5r83sq5wv8o8riZXZXVCI7OwGB5/iKH3kAhGIfK
nBvSVnEl6Wx+bZqWu23aMxUg2+Lg533aFj9aeEXC13pMrn5JMgF+C2I9mkiYo3QA
Y96AXexEFRHBwE/JV1w7FibGOtS4BvQqD8gF8XJNN7Y6n3PpsaRZ++I0P0FcUPv0
5DjGR2Q1mx1Hh1DJbnRg05AsIbBXHEOx4TjSy9b/PO7TVA5pQggTHNec09QdpXBs
mkCrR08AvKigSFlY7fGKYwu2K3qwwmeWblqW5/eWpaSU0iwVSdASZG1d7Cv7gXBd
po4f+d6aD1uDLaN9cE0a6iyO3kf+ZIdHzd4KcwIDAQABAoIBAH8ev5vZ9oFli+IG
PrfN897p7gG24aoRzxdjWqzoqsX+sh7AjKMnjeEKPyNZRKpOX/OyjVQdwAG3dGIa
wshr8Xf7NGqmY7E6972KEqIgth5PVo6GMklKr5ZItc1DHZoWkKmvUuJqO2TAIMfa
MZ5Ofe2mXxX51+2ux8palNBJESN4nbFwJw0b/R6fq5OlipHd5laaVe/QXm4ryBo6
vr/HWu71kWyF+/HsWtMutI2A0P+33AHyEcVk0vtA7kWREdL3yXmsL8Hm8jqWdP1S
bpmdC71dj0BpCYFS5yVOuAay7dJz4+panPX7KBaIdcTMPPid6dCw0ztxGh7EZYmi
Guv9KLkCgYEA+QfAmBpqU1dpp4Mw+1PW9usIl20yhqbTdX/vwxRd5VcptEpHcb1w
nnX1NpenQpIj8rstqdwEROziRKHZK8qLMTEiPLHk53hdMk6Z/9UdJHbABPvzp1AB
bQu4PpbZwDtTtY5or9h0r48cUwjvfT4j/n/abHLKh6PEOJycM67uRS8CgYEAwWFf
Pq0hF1SlnxAnRNKk9flDhEyWwZBsOtIZR0dpihTuljn7iVXfkFv1UCnNnQMh9Tol
g1z/xCcPV6SXg94GCc8xkyw+J21J8sHAGNsJFGQ6ArS/XLMbid1r0NXR4VldX1R6
4cwKdKrfbrNnGAGTAhy0lFNShwwYMjYA9dSiRf0CgYA/YprB3E9d4SzyRzErd16K
wK4SJOgsX8AI80Rqqf9wRWxHCHUA3VAR9UIx4A3houLlgIER7/9iL80z3OIzBD3D
ipcFTd5OkFNgX6NQ+8SMKHGdkyekWXfTcp01yR2pkTAwUQwSXgztNobmF6slfLCa
sZ495kXommVyZ2JWwVrCXwKBgQC+0nUIBhNnUFH2ehwl366EQqoLPQBulTMXgAcN
vTw505nzh9fcl256pyOVLQsGavbxY6Vs0TJZvyl2lKYmq8pNl7UVw0y53zBfai1C
2bFF+/j6fp1uvhboniQr+TKYKnTnAxgXBB81LQA53rJWkAceyHCxBN+k/5xIv92G
t4JBiQKBgFWaJZGznXi70MkXPgIVAaRW1saQQjN4pBCHmmHttTUP4wzTftu8amfX
dGLybeNhOfjk8Ibof7oAksG1YknecPhDYcnsB3aTegNQhBlJ03mDlhiEwq4roAqj
fugwl2R114U+R6on5tyKY3vTCamRYpo+cpaeajsPGUeHSm9e4tKV
-----END RSA PRIVATE KEY-----`
var logger = log.LoggerError

func RsaEncrypt(msg []byte, publickeypath string) (cipherText []byte, err error) {
	defer _recovr()
	var buf []byte
	if publickeypath != "" {
		if util.IsFileExist(publickeypath) {
			if buf, err = os.ReadFile(publickeypath); err != nil {
				logger.Error("rsa encrpt error:", err)
			}
		} else {
			err = errors.New(publickeypath + " not exist")
			logger.Error(publickeypath, " not exist")
		}
	}
	if publickeypath == "" || err != nil {
		buf = []byte(defaultPubKey)
	}
	pubDecodeBlock, _ := pem.Decode(buf)
	if parsePublicKey, err := x509.ParsePKIXPublicKey(pubDecodeBlock.Bytes); err == nil {
		publicKey := parsePublicKey.(*rsa.PublicKey)
		cipherText, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
	}
	return
}

func RsaDecrypt(cipherText []byte, privatekeypath string) (msg []byte, err error) {
	defer _recovr()
	var buf []byte
	if privatekeypath != "" {
		if util.IsFileExist(privatekeypath) {
			if buf, err = os.ReadFile(privatekeypath); err != nil {
				logger.Error("rsa decrpt error:", err)
			}
		} else {
			err = errors.New(privatekeypath + " not exist")
			logger.Error(privatekeypath, " not exist")
		}
	}
	if privatekeypath == "" || err != nil {
		buf = []byte(defaultPriKey)
	}
	priDecodeBlock, _ := pem.Decode(buf)
	if privateKey, err := x509.ParsePKCS1PrivateKey(priDecodeBlock.Bytes); err == nil {
		msg, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	}
	return
}

func _recovr() {
	if err := recover(); err != nil {
	}
}
