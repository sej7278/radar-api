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

func vulnsByAsset(id string, apikey string) (err error) {
	// make http request
	res, err := makeRequest("/assets/"+id, apikey, nil)
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
	var assetDetails map[string]any
	err = json.Unmarshal(resBody, &assetDetails)
	if err != nil {
		return err
	}

	// pretty print the asset details
	fmt.Printf("Asset ID:\t%v\n", assetDetails["id"])
	fmt.Printf("Host:\t\t%v (%v)\n", assetDetails["hostname"], assetDetails["ip"])
	fmt.Printf("OS:\t\t%v %v (%v)\n", assetDetails["os"], assetDetails["os_release"], assetDetails["kernel_release"])
	fmt.Printf("Radar version:\t%v\n", assetDetails["last_inspector_version"])

	// format the timestamp
	lastUploaded, _ := time.Parse(time.RFC3339, assetDetails["last_uploaded"].(string))
	fmt.Printf("Last scan:\t%v\n", lastUploaded.Local().Format(time.RFC1123))

	// pretty print the vulnerabilities
	fmt.Printf("Vulns:\t\tC=%v, H=%v, M=%v, L=%vm\n\n", assetDetails["severity_critical"], assetDetails["severity_high"], assetDetails["severity_medium"], assetDetails["severity_low"])

	// return nil if no errors
	return nil
}

func listAssets(apikey string) (ids []string, err error) {
	// make http request
	res, err := makeRequest("/assets", apikey, nil)
	if err != nil {
		return nil, err
	}

	// read the response
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal the response
	var response []any
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}

	// loop through assets and extract ids
	for _, asset := range response {
		id := fmt.Sprintf("%.0f", asset.(map[string]any)["id"].(float64))
		ids = append(ids, id)
	}

	// return the list of ids
	return ids, nil
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

	// fetch a list of asset ids
	ids, err := listAssets(apikey)
	if err != nil {
		log.Fatal(err)
	}

	// fetch vulns for each id
	for _, id := range ids {
		err = vulnsByAsset(id, apikey)
		if err != nil {
			log.Fatal(err)
		}
	}
}
