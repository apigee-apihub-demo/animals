#!/bin/sh

# Deploy service and sidecar to Cloud Run.
gcloud beta run services replace service.yaml --region us-west1

# Allow all users to call the service (ESPv2 will manage access).
gcloud run services set-iam-policy --region us-west1 zebras-example iam.yaml
