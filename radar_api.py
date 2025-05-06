#!/usr/bin/env python3

# pylint: disable=line-too-long

"""Radar api demo."""
import sys
from pathlib import Path
import requests


def main():
    """Make a request to the /assets endpoint."""
    try:
        # load apikey from ~/.radarapi
        home = Path.home()
        with open(home / ".radarapi", "r", encoding="utf-8") as f:
            apikey = f.read().strip()
    except (PermissionError, FileNotFoundError, OSError) as e:
        print(e)
        sys.exit(1)

    try:
        # make a get request
        r = requests.get(
            "https://radar.tuxcare.com/external/assets",
            headers={"X-API-KEY": apikey, "Content-Type": "application/json"},
            timeout=5,
        )

        # print any errors
        print(r.raise_for_status())

        # # loop through the json response
        for i in r.json():
            print(f"Asset ID:\t{i.get('id', '')}")
            print(f"Host:\t\t{i.get('hostname', '')} ({i.get('ip', '')})")
            print(
                f"OS:\t\t{i.get('os', '')} {i.get('os_release', '')} ({i.get('kernel_release', '')})"
            )
            print(f"Radar version:\t{i.get('last_inspector_version', '')}")
            print(f"Last scan:\t{i.get('last_uploaded', '')}\n")

    except requests.exceptions.RequestException as e:
        # handle exceptions
        print(e)


if __name__ == "__main__":
    main()
