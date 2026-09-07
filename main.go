package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func makeRequest(url string, apikey string, bodyReader io.Reader) (res *http.Response, err error) {
	// construct http request
	req, err := http.NewRequest("GET", "https://radar.tuxcare.com/external"+url, bodyReader)
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apikey)

	// make request
	var client *http.Client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	// check for http errors
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %v", res.Status)
	}

	// return response if no errors
	return res, nil
}

func assetVulns(apikey string) (err error) {
	// make http request
	res, err := makeRequest("/assets", apikey, nil)
	if err != nil {
		return err
	}

	// read the response
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshal the response
	var response []any
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return err
	}

	// loop through the assets
	for _, asset := range response {
		assetDetails, ok := asset.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid asset response")
		}

		if _, ok := assetDetails["id"].(float64); !ok {
			return fmt.Errorf("asset response is missing a valid id")
		}

		assetInfo := assetDetails

		// pretty print the asset details
		fmt.Printf("Asset ID:\t%v\n", assetInfo["id"])
		fmt.Printf("Host:\t\t%v (%v)\n", assetInfo["hostname"], assetInfo["ip"])
		fmt.Printf("OS:\t\t%v %v (%v)\n", assetInfo["os"], assetInfo["os_release"], assetInfo["kernel_release"])
		fmt.Printf("Radar version:\t%v\n", assetInfo["last_inspector_version"])

		// format the timestamp
		if lastUploadedStr, ok := assetInfo["last_uploaded"].(string); ok {
			if lastUploaded, err := time.Parse(time.RFC3339, lastUploadedStr); err == nil {
				fmt.Printf("Last scan:\t%v\n", lastUploaded.Local().Format(time.RFC1123))
			}
		}

		// pretty print the vulnerabilities
		vulnerabilities, ok := assetInfo["vulnerabilities"].(map[string]any)
		if !ok {
			return fmt.Errorf("asset response is missing vulnerabilities")
		}
		fmt.Printf("Vulns:\t\t🟣=%v, 🔴=%v, 🟠=%v, 🟢=%v\n\n", vulnerabilities["severity_critical"], vulnerabilities["severity_high"], vulnerabilities["severity_medium"], vulnerabilities["severity_low"])
	}

	// return nil if no errors
	return nil
}

func readConfig() (apikey string, err error) {
	// find home directory independent of os
	var home string
	home, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// append home directory to config file location
	var ConfigFile = filepath.Join(home, ".radarapi")

	// open config file for reading
	file, err := os.Open(ConfigFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// read the api key from the config file
	_, err = fmt.Fscanf(file, "%s", &apikey)
	if err != nil {
		return "", err
	}

	// return the api key and no error
	return apikey, nil
}

func main() {
	// read the api key from the config file
	apikey, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// fetch assets and vulns
	err = assetVulns(apikey)
	if err != nil {
		log.Fatal(err)
	}
}
