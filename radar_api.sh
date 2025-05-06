#!/bin/bash

# get api key from ~/.radarapi file
if [[ -f ~/.radarapi ]]; then
  API_KEY=$(<~/.radarapi)
else
  echo "API key file not found. Please create ~/.radarapi with your API key."
  exit 1
fi

BASE_URL="https://radar.tuxcare.com/external"

# csv header
echo "Asset ID,Hostname,IP,OS,Last Analyzed,Risk Score,Critical,High,Medium"

curl -s -H "x-api-key: $API_KEY" "$BASE_URL/assets" | jq -c '.[]' | while read -r asset; do
    id=$(echo "$asset" | jq -r '.id')
    hostname=$(echo "$asset" | jq -r '.hostname')
    ip=$(echo "$asset" | jq -r '.ip')
    os=$(echo "$asset" | jq -r '.os')
    last_analyzed=$(echo "$asset" | jq -r '.last_analyzed')

    detail=$(curl -s -H "x-api-key: $API_KEY" "$BASE_URL/assets/$id")

    risk_score=$(echo "$detail" | jq -r '.risk_score // 0')
    critical=$(echo "$detail" | jq -r '.severity_critical // 0')
    high=$(echo "$detail" | jq -r '.severity_high // 0')
    medium=$(echo "$detail" | jq -r '.severity_medium // 0')

    if [[ "$critical" -eq 0 && "$high" -eq 0 && "$medium" -eq 0 ]]; then
        vulns=$(curl -s -H "x-api-key: $API_KEY" "$BASE_URL/assets/$id/vulnerabilities")
        critical=$(echo "$vulns" | jq '[.[] | select(.severity=="critical")] | length')
        high=$(echo "$vulns" | jq '[.[] | select(.severity=="high")] | length')
        medium=$(echo "$vulns" | jq '[.[] | select(.severity=="medium")] | length')
    fi

    echo "$id,$hostname,$ip,$os,$last_analyzed,$risk_score,$critical,$high,$medium"
done
