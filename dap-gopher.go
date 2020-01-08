package main

import (

	//"strings"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/fatih/color"
)

func main() {

	color.Cyan("Retrieving the access-token from /run/conjur/access-token...")

	buf, err := ioutil.ReadFile("/run/conjur/access-token")
	if err != nil {
		log.Fatalln(err)
	}

	color.Cyan("jwt as we read it....")
	jwtString := string(buf)
	color.Magenta(jwtString)

	color.Cyan("jwt base64 encoded...")
	uEnc := b64.StdEncoding.EncodeToString([]byte(buf))
	color.Magenta(uEnc)

	// Assumes the env has been set in the container by the configMap
	var conjurApplianceURL = os.Getenv("CONJUR_APPLIANCE_URL")
	var conjurAccount = os.Getenv("CONJUR_ACCOUNT")

	// Hardcoded - should move this out into the env
	secret := "secrets/frontend/nginx_pwd"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := conjurApplianceURL + "/secrets/" + conjurAccount + "/variable/" + secret

	color.Cyan("url: " + url)

	//client := &http.Client{}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	authHeader := " Token token=\"" + uEnc + "\""
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "text/html")

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	color.Cyan(string(requestDump))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		color.Red("HTTP Response Status: %d %s.", resp.StatusCode, http.StatusText(resp.StatusCode))
	} else {
		color.Green("HTTP Response Status: %d %s.", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	result := string(body)
	color.Cyan("\nThe secret %s is: ", secret)
	color.White(result)
	color.Cyan("--------------------\n")
}

// secrets/frontend/nginx_pwd
// $(cat /run/conjur/access-token | base64 | tr -d '\r\n')
// echo "Grabbing Authentication token"
//         token=$(cat /run/conjur/access-token)
//         echo "Here's the access token:"
//         echo "$token"
//         echo "Formatting token"
//         formatted_token=$(cat /run/conjur/access-token | base64 | tr -d '\r\n')
//         echo "Formatted token is:"
//         echo "$formatted_token"
//         echo "Making call to retrieve $SECRET"
//         output=$(curl -k -s -X GET -H "Authorization: Token token=\"$formatted_token\"" $CONJUR_APPLIANCE_URL/secrets/$CONJUR_ACCOUNT/variable/$SECRET)
//         echo "Here's the output:"
//         echo "$output"
//         echo "-----"
//         sleep 5
