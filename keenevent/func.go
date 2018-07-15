package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	fdk "github.com/fnproject/fdk-go"
)

// Event models an event logged to Keen.io
type Event struct {
	Trigger string `json:"trigger"`
	Event   string `json:"event"`
	Payload map[string]interface{}
}

// Output models messages to be displayed in the function response
type Output struct {
	Log []string
}

// HandleErr models errors that may occur during function startup
type HandleErr struct {
	Log string
}

// Serve satisfies the fdk Handler interface
func (herr HandleErr) Serve(ctx context.Context, in io.Reader, out io.Writer) {
	var outmsg = map[string]string{"log": ""}

	outmsg["log"] = herr.Log
	if len(outmsg["log"]) == 0 {
		// no message supplied, report a generic message
		outmsg["log"] = "an error occurred"
	}

	json.NewEncoder(out).Encode(&outmsg)
}

const envVar = "FUNC_URL"

var (
	err     error
	ok      bool
	keenURL string
	output  Output
)

func main() {
	// Get the Keen.io URL from an environment variable
	keenURL, ok = os.LookupEnv(envVar)
	if !ok {
		// missing environment variable, function call with error message
		log.Println("ERROR: missing url environment variable")
		fdk.Handle(HandleErr{Log: "error: missing url environment variable"})
		return
	}

	// Validate the Keen URL
	_, err = url.ParseRequestURI(keenURL)
	if err != nil {
		// invalid Keen.io url specified
		log.Println(fmt.Sprintf("ERROR: invalid Keen.io url; see - %v", err))
		fdk.Handle(HandleErr{Log: fmt.Sprintf("error: invalid Keen.io url; see - %v", err)})
		return
	}

	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// Create an event object to be logged to Keen.io
	ev := &Event{Trigger: "fn", Event: "log"}

	// Grab the inbound json request data
	req := map[string]interface{}{}
	err = json.NewDecoder(in).Decode(&req)
	if err != nil {
		// error occurred decoding the inbound request
		log.Printf("ERROR: decoding the request payload; no event data sent to Keen.io")
		output.Log = append(output.Log, "error decoding the request payload; no event data sent to Keen.io")
		json.NewEncoder(out).Encode(&output)
		return
	}
	ev.Payload = req

	// Log the Keen.io event
	evBytes, err := json.Marshal(ev)
	if err != nil {
		// error marshalling the Keen event data in a byte slice
		log.Printf("ERROR: error marshalling the Keen event data in a byte slice; no event data sent to Keen.io")
		output.Log = append(output.Log, "error marshalling the Keen event data in a byte slice; no event data sent to Keen.io")
		json.NewEncoder(out).Encode(&output)
		return
	}
	evRdr := bytes.NewReader(evBytes)

	// Log the event and payload to Keen.io
	resp, err := http.Post(keenURL, "application/json", evRdr)
	if err != nil {
		// error occurred logging an event to Keen.io
		log.Printf("ERROR: error occurred logging the event\n***\n%v\n***\n to Keen.io. See: %v\n",
			string(evBytes),
			err)
		output.Log = append(output.Log, fmt.Sprintf("error occurred logging the event to Keen.io. See: %v", err))
		json.NewEncoder(out).Encode(&output)
		return
	}

	// Keen event logged
	output.Log = append(output.Log, fmt.Sprintf("event logged to Keen.io with status: %v", resp.Status))
	json.NewEncoder(out).Encode(&output)
}
