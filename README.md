# Fan Control
A Go application to control the fan speed based on temperature on iDRAC 6 servers

## Build
Clone the repo and run go build:
`go build .`

## Usage 
Edit the `config.yaml` file with the hostname, username, and password for IPMI.

To get current temperature:
`./fan_control get temp`

To run the program continuously:
`./fan_control run`
