# transcribe
Practice Go by building a cloud native application to transcribe audio to text using a google cloud function

## Local Development

To run the google function-framework-go you start the function by:

```
LOCAL_ONLY=true FUNCTION_TARGET="TranscribeCloudFunc" go run cmd/main.go
```

To call the function curl to localhost with a cloud event as the payload:

```
curl -X POST http://localhost:8080 \
-H "Content-Type: application/json" \
-H "Ce-Specversion: 1.0" \
-H "Ce-Type: your.event.type" \
-H "Ce-Source: your.event.source" \
-H "Ce-Id: your-event-id" \
-d '@./cloudevents/StorageObjectDataWorkflows.json'

```

## Deploying Cloud Function

Enable Gen2 APIs (on-time)

```
gcloud services enable \
eventarc.googleapis.com \
pubsub.googleapis.com \
run.googleapis.com \
storage.googleapis.com
```

### Grant roles for Events

See https://cloud.google.com/eventarc/docs/run/create-trigger-storage-gcloud#local-shell

```
gcloud projects add-iam-policy-binding digest-427612 \
    --member=serviceAccount:672553277717-compute@developer.gserviceaccount.com \
    --role=roles/eventarc.eventReceiver
```

```
SERVICE_ACCOUNT="$(gsutil kms serviceaccount -p digest-427612)"

gcloud projects add-iam-policy-binding digest-427612 \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role='roles/pubsub.publisher'
```

Deploy Function

```
gcloud functions deploy TranscribeCloudFunc --entry-point TranscribeCloudFunc --runtime go121 --trigger-bucket digest-bucket --allow-unauthenticated --region northamerica-northeast1 --gen2
```

