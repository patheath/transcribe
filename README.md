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