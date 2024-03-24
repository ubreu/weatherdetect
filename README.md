# RainDetect

This is a cloud function which determines whether precipitation is detected for any of the given station.

## Development

Test the cloud function locally:
````
export FUNCTION_TARGET=DetectForLocation
go run cmd/main.go
````