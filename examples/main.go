package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"

	libdnsnextdns "github.com/evansmurithi/libdns-nextdns"
	"github.com/libdns/libdns"
)

func main() {
	ctx := context.Background()
	profileID := os.Getenv("NEXTDNS_PROFILE_ID")

	p, err := libdnsnextdns.NewProvider(
		ctx,
		libdnsnextdns.Opt{
			ApiKey: os.Getenv("NEXTDNS_API_KEY"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	records, err := p.GetRecords(ctx, profileID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", records)

	// Generate a random IP address
	ip := net.IPv4(
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	).String()

	recordsToCreate := []libdns.Record{
		{
			Name:  "test.example.com",
			Value: ip,
			Type:  "A",
		},
	}

	createdRecords, err := p.AppendRecords(ctx, profileID, recordsToCreate)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("\nCreated records: %+v\n", createdRecords)

	deletedRecords, err := p.DeleteRecords(ctx, profileID, createdRecords)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("\nDeleted records: %+v\n", deletedRecords)
}
