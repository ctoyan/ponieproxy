package ponieproxy

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDVzCCAj+gAwIBAgIUNGnaeDvVF1PM3bFU6QXp0c9o6LEwDQYJKoZIhvcNAQEL
BQAwOzELMAkGA1UEBhMCVUsxDzANBgNVBAcMBkxvbmRvbjEbMBkGA1UECgwSSW50
IHRlc3RzIFByb3h5IENBMB4XDTIwMDIyMDEyNDA0OVoXDTQwMDIxNTEyNDA0OVow
OzELMAkGA1UEBhMCVUsxDzANBgNVBAcMBkxvbmRvbjEbMBkGA1UECgwSSW50IHRl
c3RzIFByb3h5IENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0pVw
fJ3uP5omHvq4GwycqRTL+2IvyTgLJhr9RxJ3Gl0SJ0K1Drxgtx5bYk3s4PE35ign
KCLa2rP0FLC0wOyRth1x1ZR0kjEvG2rB/rDjo0lgdwpvaDjBrTrfzdH8GKcYGnti
Z5iPaxpA6h654uyms7/GzN+8AlegBeLRJPo4l5kHHgtgV5dZDSnPWOsWnPFJgIA8
7NkyW2jGI/1hO7+vLBjUijbgdqutmF6YUl7Ui9BFjBZk+8W9iHIA27z6rxb/1uFr
yJ7waWo/mPJg/3pAThJ7SXH/ZNS52r12Wu+d1H2mIrNH2QZQaTymF+Xf2jtbeQ/N
dwQVxZeEsbvLu+aeAwIDAQABo1MwUTAdBgNVHQ4EFgQUem0U5X4yj2Vxg0ghAv/Z
VYArSkswHwYDVR0jBBgwFoAUem0U5X4yj2Vxg0ghAv/ZVYArSkswDwYDVR0TAQH/
BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEATIm6gdf2POtM+zA/oPpxXaFJOFPu
bPfgfGAgbU2C0SeLp+DX27J15ULMzcb7gHauHFfK7O9OiQLtqZM6e0DReajl9R9q
b4NM3BXq+QzjpBnz3nmszTrb0sJUVFjFD80la9V05kbdjxrjJK4fIrGhs8BA82zw
yOfQOHT81v1P0+ieihM0+FcCCu3vJJqTEmYutLbo4/FVSbxZcUrV9NPRp1dn2zX+
CA1CKtgVdD+OSpqG2HXvvuEQi7ydvQo8svgFS4ZRkzn/yG2DHovl9ObeJH49sbJM
PeW7VKXrVNYnWtfnBwnL3vzHDtmJXiLAseuQZjeaJOAwy9lLvPKFnthYPA==
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDSlXB8ne4/miYe
+rgbDJypFMv7Yi/JOAsmGv1HEncaXRInQrUOvGC3HltiTezg8TfmKCcoItras/QU
sLTA7JG2HXHVlHSSMS8basH+sOOjSWB3Cm9oOMGtOt/N0fwYpxgae2JnmI9rGkDq
Hrni7Kazv8bM37wCV6AF4tEk+jiXmQceC2BXl1kNKc9Y6xac8UmAgDzs2TJbaMYj
/WE7v68sGNSKNuB2q62YXphSXtSL0EWMFmT7xb2IcgDbvPqvFv/W4WvInvBpaj+Y
8mD/ekBOEntJcf9k1LnavXZa753UfaYis0fZBlBpPKYX5d/aO1t5D813BBXFl4Sx
u8u75p4DAgMBAAECggEBAM6LdMFFxiDr+Of50gn13NKaa1gtfeFG7uh5IGNNYTSi
hOOtMhk5+0Kgq9FHzUb0UXeVepKLXU2Vo5mHmEKBxolxZ+2spomdZC7oD07YAO5v
UgZwXcVDpfNbA5jElRc5DRMsYeBqaoEKRxFbAcfphbhYKY1ZBPbnKzWaurgiFM/a
qB64BP24slm3zqId7YhblYLrv6BGK2j2FmB4sIs29zpSQ1g/6muroyA6pZdcd6Ht
BRLR9FxWI7qSQbcZnn3u+hStEJhG9cIuZVnOBmG4+hh4T0d8qg9nVr4BHKXebYaR
oBGpn5plYcghR3xfeX8P8WHu5jLprLAlWeO+JqjELIECgYEA9xdIEWod3qKPSNA7
EOj1FwD2g3Ov7vmtAfRKgmhnDnlLs+nlezjvE0rEL6OjTdSUfUcwp3LHSdk/Jlfm
nVlkmHPQVRjXAkqahUNt9j7ZnRJxDh/9u+uVl3rguwqV3aKiO+NyV2b+7Hh+Ir8J
Yzx8aBmXJ7mNCO5xa0IbRFTew/kCgYEA2i0vru7CFpZnQqnFyNMDT4d3XNamJGKG
Jcf22157wp/3OX/WyQOpnl2rdoiuyTcaffBgBk/DKLapfEQSLTH1mQsKDuxZinzY
7/yMimxDXc2K5ubyAuA/E/PCf/YsPYO121bS4z03SsfKhquDB9pWy+wKBuA5uFa+
EP0Wh/LKuNsCgYBJHdGMnasjE1V1BXFFCrpjyTwpH9Wi0K0aU/CscDp2tPvqzD7E
3M8aFVjChBix0kLyY1uJYVSJjMi8DuzGCQrUdgji9YvCONNKte5XHLgGW8uqk1rg
/dBxV8Iidvpr8FEziZVvOaIb1Xf1zjP38pEZuODat3R9fRmA1Ln+2WJl+QKBgGf8
iXGTEqa5YNYBHOeuyzEom9d/5wgIfW+ccyfzTIFixO5+49xDBqEYfBSu6L+2p8XG
v73CXn4VMYqs1wz7dtdOz6h1NegvwSYA9Os01pbq1H1hLY/5WZck41sh9cwL7q1w
IGt2TdgyiXDOZlFj22KuScklLd2SWly2g/qf2cdpAoGBAO7uhjJn2j5aduqB/PFg
2zFS4X52Ikk/5gQomcPlKy6ubhtG9uE5UoEc6ftWdBZZ1MR7QRMw91ToszcWYS0b
CGKiKzvpNX71BF+JfIRAm4mPXIqRQJj8nlNfplWVBLKoVqLI52p/PIV+SCp0iN/g
v/ZO4iGry0PlqPoFt/ItaDI7
-----END PRIVATE KEY-----`)

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
