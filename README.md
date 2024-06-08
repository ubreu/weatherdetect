# WeatherDetect

This is a cloud function which determines whether certain weather features/components (i.e. sunshine or precipitation) are detected for any of the given stations.

Detect for a list of the given stations:
````
curl https://europe-west6-weather-detect.cloudfunctions.net/weather-detect?station=WAG&station=SMA&station=REH&station=DIT&station=ZWK
````

## Development

Test the cloud function locally:
````
export FUNCTION_TARGET=DetectForLocation
go run cmd/main.go
curl "localhost:8080?station=SMA&station=BRZ&station=KLO&station=WFJ&station=BRT&station=GOS"
````


## Deployment

Initialize the gcloud CLI (only run once)
````
gcloud init
````

Authorize:
````
gcloud auth login
````

Set project:
````
gcloud config set project weather-detect
````

Deploy:
````
gcloud functions deploy weather-detect \
--gen2 \
--runtime=go122 \
--region=europe-west6 \
--source=. \
--entry-point=DetectForLocation \
--trigger-http \
--allow-unauthenticated
````