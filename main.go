package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DNSResponse struct {
	Result []DNSRecord `json:"result"`
}

func main() {
	apiToken := os.Getenv("API_TOKEN")
	zoneID := os.Getenv("ZONE_ID")
	recordName := os.Getenv("RECORD_NAME")

	if apiToken == "" || zoneID == "" || recordName == "" {
		fmt.Println("Error: Missing required environment variables.")
		os.Exit(1)
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	resp, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		log.Fatal(err)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ownIP := string(bytes.TrimSpace(ip))

	recs, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Fatal(err)
	}

	var recordToUpdate *cloudflare.DNSRecord
	for _, r := range recs {
		if r.Name == recordName {
			recordToUpdate = &r
		}
	}

	if recordToUpdate == nil {
		log.Fatal("Error: DNS record not found.")
	}

	if recordToUpdate.Content == ownIP {
		fmt.Println("No update needed: IP address has not changed.")
		return
	}

	rec, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID),
		cloudflare.UpdateDNSRecordParams{ID: recordToUpdate.ID, Type: recordToUpdate.Type, Name: recordToUpdate.Name, Content: ownIP})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully updated IP address to: %s\n", rec.Content)
}
