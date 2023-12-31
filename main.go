package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	hostname := os.Getenv("DDNS_HOSTNAME")
	if hostname == "" {
		log.Fatal("Failed to get DDNS_HOSTNAME from environment.")
	}

	zone := os.Getenv("DDNS_ZONE")
	if zone == "" {
		log.Fatal("Failed to get DDNS_ZONE from environment.")
	}

	ip := os.Getenv("DDNS_IP")
	if ip == "" {
		log.Fatal("Failed to get DDNS_IP from environment.")
	}

	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Failed to get CLOUDFLARE_API_TOKEN from environment.")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		log.Fatal("Failed to get Zone ID.  ", err)
	}

	zoneRC := cloudflare.ZoneIdentifier(zoneID)

	var listParams cloudflare.ListDNSRecordsParams
	listParams.Type = "A"
	listParams.Name = hostname

	records, _, err := api.ListDNSRecords(ctx, zoneRC, listParams)
	if err != nil {
		log.Fatal("Failed to list DNS records.  ", err)
	}

	if len(records) != 1 {
		log.Fatal("Did not get exactly one matching A record.  ", records)
	}

	record := records[0]
	if record.Content == ip {
		os.Exit(0)
	}

	var updateParams cloudflare.UpdateDNSRecordParams
	updateParams.ID = record.ID
	updateParams.Name = record.Name
	updateParams.Content = ip

	_, err = api.UpdateDNSRecord(ctx, zoneRC, updateParams)
	if err != nil {
		log.Fatal("Failed to update DNS record.  ", err)
	}

	log.Print("Updated DNS record: ", hostname, ", IP: ", record.Content, " -> ", updateParams.Content)
	os.Exit(0)
}
