package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

// geturls constructs a slice of web site addresses from enviroment variables
func geturls() ([]string, error) {
	var (
		err             error
		ok              bool
		curvar, curval  string
		sites, errSites []string
	)
	const prefx = "FUNC_SITE_"

	// Iterate through a maximum number of environment variables
	for i := 0; i < 10; i++ {
		curvar = fmt.Sprintf("%v%v", prefx, i)

		// Get a site from an environment variable
		curval, ok = os.LookupEnv(curvar)
		if !ok {
			// no more variables, stop
			break
		}

		// Validate the URL
		_, err = url.ParseRequestURI(curval)
		if err != nil {
			// URL in error
			errSites = append(errSites, curval)
		}

		sites = append(sites, curval)
	}

	// Any errors?
	err = nil
	if len(errSites) != 0 {
		// one or more URLs are in error
		err = fmt.Errorf("one or more URLs in error: %v", strings.Join(errSites, ", "))
	}

	if len(sites) == 0 {
		// no valid URLS found
		err = fmt.Errorf("no valid URLs found %v", err)
	}

	return sites, err
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	var (
		err    error
		sites  []string
		output []string
		resp   *http.Response
	)

	// Set up to write to 'out'
	toOut := bufio.NewWriter(out)

	// Get the list of web sites to check
	sites, err = geturls()
	if err != nil {
		toOut.WriteString(fmt.Sprintf("ERROR: %v", err))
		return
	}

	// Iterate through the sites and "GET" each one
	for i := 0; i < len(sites); i++ {
		resp, err = http.Get(sites[i])

		if resp.StatusCode != 200 || err != nil {
			// not ok response; report as an error
			output = append(output, fmt.Sprintf("SITE: %v - RESULT: Error - %v", sites[i], err))
			continue
		}

		// ok response
		output = append(output, fmt.Sprintf("SITE: %v - RESULT: Ok - %v", sites[i], resp.Status))
	}

	// Write results
	toOut.WriteString(fmt.Sprintf("FUNCTION OUTPUT:\n\n%v", strings.Join(output, "\n")))
}
