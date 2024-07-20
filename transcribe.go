package trascribe

import (
	"context"
	"log"

	speech "cloud.google.com/go/speech/apiv2"
	speechpb "cloud.google.com/go/speech/apiv2/speechpb"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// We only care about payload in the cloud event
// Tried to use the cloudevents sdk but it did not like the 
// spec version of the cloud event (the Test Trigger in console)
// sets the spec version as a float (1.0) and cloudevents sdk expects
// it to be a string
type MyCloudEvent struct {
    Data MyStorageObjectData `json:"data"`
    Id string `json:"id"`
}

// Similarly could not use storagedata "github.com/googleapis/google-cloudevents-go/cloud/storagedata"
// as it expects the metageneration to be a string but was getting a float
// So had to recreate the struct with what I needed only, what a pain
type MyStorageObjectData struct {
    Bucket string `json:"bucket"`
    Name string `json:"name"`
    ContentType string `json:"contentType"`
    LocalOnlyTest bool `json:"localOnlyTest"`  // Custom field for testing, won't call speech to text if true
}

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("TranscribeCloudFunc", transcribeCloudFunc)
}

// Function TranscribeCloudFunc accepts and handles a CloudEvent
// when a new audio file is uploaded to a Cloud Storage bucket.
func transcribeCloudFunc(ctx context.Context, e event.Event) error {
	// log.Println(e)

    event := MyCloudEvent{}
    if err := e.DataAs(&event); err != nil {
        log.Fatal("e.DataAs: ", err)
    }
 
    // Check if the content type is audio
    if event.Data.ContentType != "audio/x-wav" {
        log.Printf("Exiting - content type is not audio: %v", event.Data.ContentType)
        return nil
    }
    // Check if the folder is the right one
    if event.Data.Name[:12] != "audio-files/" {
        log.Printf("Exiting - folder is not audio-files: %v", event.Data.Name[:12])
        return nil
    }

    log.Printf("We have file: %v to transcribe to text.  Starting....", event.Data.Name)
    if !event.Data.LocalOnlyTest {
        toText(event.Data.Name)
        log.Println("Transcription has been started.")
    } else {
        log.Println("LocalOnlyTest is true, will not call speech to text.")
    }
	return nil
}

func toText(name string) {
	ctx := context.Background()
	c, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// The path to the remote audio file to transcribe
	fileUri := "gs://digest-bucket/" + name
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
			EnableWordTimeOffsets: false,
			EnableWordConfidence:  false,
            EnableAutomaticPunctuation: true,
            // Although documentation says its supported Speaker Diarization is not supported in v2
            // https://www.googlecloudcommunity.com/gc/AI-ML/Speaker-Diarization-is-disabled-even-for-supported-languages-in/m-p/616388
            // DiarizationConfig: &speechpb.SpeakerDiarizationConfig{
            //    MinSpeakerCount: 1,
            //    MaxSpeakerCount: 6,
            // },
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
    
    
    var op *speech.BatchRecognizeOperation
	op, err = c.BatchRecognize(ctx, req)
	if err != nil {
		log.Fatalf("failed to create BatchRecognize: %v", err)
	}
    log.Printf("Created Speech to Text operation: %s", op.Name())
}
