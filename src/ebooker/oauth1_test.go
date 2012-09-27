package ebooker

import (
	"launchpad.net/gocheck"

	"fmt"
	"net/http"
	"strings"
)

// hook up gocheck into the gotest runner.
type OAuthSuite struct{}

var _ = gocheck.Suite(&OAuthSuite{})

func (o OAuthSuite) TestPercentEncode(c *gocheck.C) {

	testCases := map[string]string{
		"hello":                          "hello",
		"CAPS":                           "CAPS",
		"11WithDigits9":                  "11WithDigits9",
		"a space and exclamation point!": "a%20space%20and%20exclamation%20point%21",
		"Dogs, Cats & Mice":              "Dogs%2C%20Cats%20%26%20Mice",
		"Reserved Chars -._~":            "Reserved%20Chars%20-._~",
		"Ladies + Gentlemen":             "Ladies%20%2B%20Gentlemen"}

	for k, v := range testCases {
		c.Assert(percentEncode(k), gocheck.Equals, v)
	}
}

// Circumventing the 'createOAuthRequest' API, this tests against the example
// Twitter themselves walk you through, once you've obtained an access token.
//
// https://dev.twitter.com/docs/auth/creating-signature
func (o OAuthSuite) TestTwitterSignatureExample(c *gocheck.C) {
    status := "Hello Ladies + Gentlemen, a signed OAuth request!"

	paramMap := map[string]string{
		"status":                 status,
		"include_entities":       "true",
		"oauth_consumer_key":     "xvz1evFS4wEEPTGEFPHBog",
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		"oauth_version":          "1.0" }

    consumerSecret := "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw"
    tokenSecret := "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"

    url := "https://api.twitter.com/1/statuses/update.json"
	req, err := http.NewRequest("POST", url, strings.NewReader(percentEncode(status)))
	if err != nil {
		fmt.Printf("error %v in making the request to %v\n", err, url)
	}

    oauthObject := OAuthRequest{paramMap, url, req, consumerSecret, tokenSecret}
    c.Assert(oauthObject.makeSigningKey(), gocheck.Equals, "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw&LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE")

    expected := "POST&https%3A%2F%2Fapi.twitter.com%2F1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520%252B%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521"
    c.Assert(oauthObject.makeSignatureBaseString(), gocheck.Equals, expected)

    oauthObject.createSignature()
    c.Assert(oauthObject.parameterStringMap["oauth_signature"], gocheck.Equals, "tnnArxj06cWHq44gCs1OSKk/jLY=")
}
