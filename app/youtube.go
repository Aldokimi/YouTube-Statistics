package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func getChannelID(channelName string, apiKey string) (string, error) {
    client := &http.Client{
        Transport: &transport.APIKey{Key: apiKey},
    }
    service, err := youtube.New(client)
    if err != nil {
        return "Couldn't create service!", err
    }
	
    call := service.Channels.List([]string{"id"}).
        ForUsername(channelName).
        MaxResults(1)
    response, err := call.Do()
    if err != nil {
        return "Couldn't get response!", err
    }

    if len(response.Items) == 0 {
        return "", fmt.Errorf("no channel found with name %s", channelName)
    }

    return response.Items[0].Id, nil
}

func getChannelsStatus(k string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()

		yts, err := youtube.NewService(ctx, option.WithAPIKey(k)); if err != nil {
			fmt.Println("failed to create service")
			w.WriteHeader(http.StatusBadRequest)
			return 
		}

		call := yts.Channels.List([]string{"id, snippet, contentDetails, statistics"})
		id, err := getChannelID("GitHub", k); if err != nil {
			fmt.Printf("Failed to get the youtube channel id: %v", err)
		}

		/*"UClSv7tWDA4wkCTLhZl1YBlw"*/
		response, err := call.Id(id).Do(); if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return 
		}

		channelData := make(map[string]interface{})

		// Store the channel ID, snippet, and statistics in the map
		channelData["id"] = response.Items[0].Id
		channelData["contentDetails"] = response.Items[0].ContentDetails
		channelData["snippet"] = response.Items[0].Snippet
		channelData["statistics"] = response.Items[0].Statistics

		// Print the channel data
		fmt.Printf("Channel data: %v", channelData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(channelData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}