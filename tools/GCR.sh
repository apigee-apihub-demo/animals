#!/bin/sh

gcloud auth configure-docker us-west1-docker.pkg.dev

docker build . --tag us-west1-docker.pkg.dev/apigee-apihub-demo/animals/animals:latest
docker push us-west1-docker.pkg.dev/apigee-apihub-demo/animals/animals:latest
