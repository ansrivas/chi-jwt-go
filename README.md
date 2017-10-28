chi-jwt:
---


Using jwt-go with go-chi:

Install:
---

* Clone the project and run

	`go get -u github.com/ansrivas/chi-jwt-go`


* Create a directory named `keys` and copy the keys from the repository

* Execute your program as `chi-jwt-go -keyPath keys`

* Install [httpie](https://github.com/jakubroztocil/httpie) ( super awesome curl alternative )

* Authenticate to the service:
```
$ http POST localhost:8080/login username="someone" password="p@assword"

HTTP/1.1 200 OK
Content-Length: 439
Content-Type: text/plain; charset=utf-8
Date: Sat, 28 Oct 2017 01:12:20 GMT
Strict-Transport-Security: max-age=63072000; includeSubDomains

{"token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDkxNTY3NDAsImlhdCI6MTUwOTE1MzE0MH0.OzHtWc2Ag1uYvmX19-HZYbBLc8_A0xdUrO9OEYof8Mrd0eUKZYMutppiXCZBShpIP0TDpHpUooYiL1qgejI4lTkHGIUxmHvFuFdc0_8vWQWG_e9HikDyy0xte3iKjVALq3IYRql65pnAHfDFNH7kNT-mwxg38INu2ps4RDPiiXsIY6puRd-NBcP5PcQRYRPOL0HqnaEzblOGoWLmMN-PgsitGa8VzJMILc-TK5nR0K5cllrlWD_7VcIW85DLyBm2JbTK5xfAEepq3o97qLEWMtG40sZXeAch68gpIehXPynLMHO8G-axbnknbsljrSnwpvhZRe8RTI2GHFTWS8HAKw"}
```

* Use `token` from previous output after `Bearer <token_here>` in below example.
```
$ http GET localhost:8080/resource Authorization:"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDkxNTY3NDAsImlhdCI6MTUwOTE1MzE0MH0.OzHtWc2Ag1uYvmX19-HZYbBLc8_A0xdUrO9OEYof8Mrd0eUKZYMutppiXCZBShpIP0TDpHpUooYiL1qgejI4lTkHGIUxmHvFuFdc0_8vWQWG_e9HikDyy0xte3iKjVALq3IYRql65pnAHfDFNH7kNT-mwxg38INu2ps4RDPiiXsIY6puRd-NBcP5PcQRYRPOL0HqnaEzblOGoWLmMN-PgsitGa8VzJMILc-TK5nR0K5cllrlWD_7VcIW85DLyBm2JbTK5xfAEepq3o97qLEWMtG40sZXeAch68gpIehXPynLMHO8G-axbnknbsljrSnwpvhZRe8RTI2GHFTWS8HAKw"

HTTP/1.1 200 OK
Content-Length: 46
Content-Type: text/plain; charset=utf-8
Date: Sat, 28 Oct 2017 01:16:55 GMT
Strict-Transport-Security: max-age=63072000; includeSubDomains

{"data":"Gained access to protected resource"}
```


References:
---
1. http://www.giantflyingsaucer.com/blog/?p=5994
2. https://gist.github.com/kiyor/7817632
3. https://github.com/ianmcmahon/encoding_ssh
4. https://github.com/XOfSpades/authentication/blob/69bdcf4131c38bfbe31b26be120ad95f4816a5ae/README.md


Create keys
------


### .p12 format keys

### Private key

`keytool -genkeypair -keystore jwtsig-test-prv-ks.p12 -storetype pkcs12 -alias jwtsigtest -keyalg RSA -keysize 2048 -sigalg SHA384withRSA -dname "CN=jwtsigtest,OU=Auth Test,O=private purpose,L=Cologne,ST=NRW,C=DE" -validity 3652`

### Public key

`keytool -exportcert -alias jwtsigtest -file jwtsig-test-pub.cert -storetype pkcs12 -keystore jwtsig-test-prv-ks.p12 -rfc`

`keytool -importcert -alias jwtsigtest -file jwtsig-test-pub.cert -storetype pkcs12 -keystore jwtsig-test-pub-ks.p12`

`rm jwtsig-test-pub.cert`

### Convert to .pem format from p12 format, this is what we will use.

### Private key

`openssl pkcs12 -in jwtsig-test-prv-ks.p12 -nocerts -out jwtsig-test-prv-ks.pem -nodes`

### Public key

Generate certificate:

`openssl pkcs12 -in jwtsig-test-pub-ks.p12 -out jwtsig-test-pub-cert.pem`

Determine public key from certificate file:

`openssl x509 -in jwtsig-test-pub-cert.pem -pubkey -noout > jwtsig-test-pub-ks.pem`
