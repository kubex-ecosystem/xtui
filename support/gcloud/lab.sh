#!/usr/bin/env bash

gcloud builds submit

gcloud container images add-tag

gcloud run deploy

gcloud config list

gcloud run services list

gcloud compute instances list

gcloud iam service-accounts list

gcloud logging read "resource.type=container "

gcloud sql instances list

gcloud run services update run-name \
  --execution-environment gen2

gcloud iam service-accounts add-iam-policy-binding [EMAIL_ADDRESS] \
