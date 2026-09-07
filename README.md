# TuxCare Radar API Demos

A series of scripts that demonstrate how to use the [Radar API](https://radar.tuxcare.com/external/docs/)

## Setting up credentials

All of the demos use the same config file. You should create a ~/.radarapi file using the following command (substitute in your API key). Note that the file is not encrypted, so we protect with UNIX permissions:

```bash
echo "b3d1115c-cce6-41e6-a0a6-f694b4701388" > ~/.radarapi
chmod 0600 ~/.radarapi
```

## Demo 1: main.go

A golang (Go) module that pretty prints to the console.

### Running from go

You can run the source script directly from go, for example:

```bash
go run main.go
```

Example output:

```text
Asset ID:       23132
Host:           ubuntu16 (192.168.0.90)
OS:             ubuntu 16.04 (4.4.0-285-generic)
Radar version:  1.6.3-1
Last scan:      Fri, 04 Sep 2026 13:00:02 BST
Vulns:          🟣=5, 🔴=17, 🟠=59, 🟢=31

Asset ID:       23202
Host:           debian11 (192.168.0.75)
OS:             debian 11 (5.10.0-46-amd64)
Radar version:  1.6.2-1
Last scan:      Mon, 17 Aug 2026 12:08:41 BST
Vulns:          🟣=0, 🔴=0, 🟠=0, 🟢=0
```

### Installing a binary

You can compile and install your platform's binary (there are no dependencies) into `$GOPATH/bin/` using:

```bash
go install github.com/sej7278/radar-api@latest
```

## Demo 2: radar_api.py

A python 3 script that prints to the console.

There is a generic shebang in the file, so you should be able to make it executable and then run it like so:

```bash
chmod u+x radar_api.py
./radar_api.py
```

Alternatively run using a specific or OS default python interpreter:

```bash
python radar_api.py
```

Example output:

```text
Asset ID:       23132
Host:           ubuntu16 (192.168.0.90)
OS:             ubuntu 16.04 (4.4.0-272-tuxcare.els43-generic)
Radar version:  1.3.0-1
Last scan:      2025-04-03T10:41:14Z

Asset ID:       23202
Host:           debian11 (192.168.0.75)
OS:             debian 11 (5.10.0-34-amd64)
Radar version:  1.4.0-beta1
Last scan:      2025-04-02T15:31:39Z
```

## Demo 3: radar_api.sh

A shell script that uses `curl` and `jq` to output CSV data.

Run it using:

```bash
chmod u+x radar_api.sh
./radar_api.sh
```

Example output (could be redirected to a file):

```text
Asset ID,Hostname,IP,OS,Last Analyzed,Risk Score,Critical,High,Medium
23132,ubuntu16,192.168.0.90,ubuntu,2025-05-06T11:03:44Z,56.0,11,36,32
23202,debian11,192.168.0.75,debian,2025-05-06T08:36:02Z,0.0,0,0,0
```
