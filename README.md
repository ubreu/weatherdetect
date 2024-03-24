# RainDetect

This is a cloud function which determines whether precipitation is detected for any of the given station.

## Detect Rain

Detect if any of the given stations is measuring precipitation:
````
curl "localhost:8080?station=SMA&station=BRZ&station=KLO&station=WFJ&station=BRT&station=GOS"
````

## Development

Test the cloud function locally:
````
export FUNCTION_TARGET=DetectForLocation
go run cmd/main.go
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

Deploy:
````
gcloud functions deploy rain-detect \
--gen2 \
--runtime=go122 \
--region=europe-west6 \
--source=. \
--entry-point=DetectForLocation \
--trigger-http \
--allow-unauthenticated
````