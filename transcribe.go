package trascribe

import (
	"context"
	"log"

	speech "cloud.google.com/go/speech/apiv2"
	speechpb "cloud.google.com/go/speech/apiv2/speechpb"
    storagedata "github.com/googleapis/google-cloudevents-go/cloud/storagedata"


	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("TranscribeCloudFunc", transcribeCloudFunc)
}

// Function TranscribeCloudFunc accepts and handles a CloudEvent
// when a new audio file is uploaded to a Cloud Storage bucket.
func transcribeCloudFunc(ctx context.Context, e event.Event) error {
	log.Println(e)

	// Your code here
	// Access the CloudEvent data payload via e.Data() or e.DataAs(...)
    // Convert event.Event DataEncoded byte array to an object
    var data storagedata.StorageObjectData

    if err := e.DataAs(&data); err != nil {
        log.Printf("failed to convert data: %v", err)
        return err
    }

    // Use the converted object
    // ...
    log.Printf("Bucket: %s, Object: %s", data.Bucket, data.Name)

	// Return nil if no error occurred
	return nil
}

func toText() {
	ctx := context.Background()
	c, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// The path to the remote audio file to transcribe
	fileUri := "gs://digest-bucket/audio-files/zoom.wav"
	// The path to the transcript result folder.
	outputUri := "gs://digest-bucket/transcripts"
	// Recognizer resource name.
	recognizer := "projects/digest-427612/locations/global/recognizers/_"

	config := &speechpb.RecognitionConfig{
		DecodingConfig: &speechpb.RecognitionConfig_ExplicitDecodingConfig{
			ExplicitDecodingConfig: &speechpb.ExplicitDecodingConfig{
				Encoding:          speechpb.ExplicitDecodingConfig_LINEAR16,
				SampleRateHertz:   32000,
				AudioChannelCount: 1,
			},
		},
		Model:         "long",
		LanguageCodes: []string{"en-US"},
		Features: &speechpb.RecognitionFeatures{
			EnableWordTimeOffsets: true,
			EnableWordConfidence:  true,
		},
	}

	audioFiles := []*speechpb.BatchRecognizeFileMetadata{
		&speechpb.BatchRecognizeFileMetadata{
			AudioSource: &speechpb.BatchRecognizeFileMetadata_Uri{
				Uri: fileUri,
			},
		},
	}
	outputConfig := &speechpb.RecognitionOutputConfig{
		Output: &speechpb.RecognitionOutputConfig_GcsOutputConfig{
			GcsOutputConfig: &speechpb.GcsOutputConfig{
				Uri: outputUri,
			},
		},
	}
	req := &speechpb.BatchRecognizeRequest{
		Recognizer:              recognizer,
		Config:                  config,
		Files:                   audioFiles,
		RecognitionOutputConfig: outputConfig,
	}
	_, err = c.BatchRecognize(ctx, req)
	if err != nil {
		log.Fatalf("failed to create BatchRecognize: %v", err)
	}
}
