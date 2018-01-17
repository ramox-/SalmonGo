package main

import (
	"encoding/json"
	"fmt"
	"os"
	"salmongo"
)

func main() {
	// create a new client pointing to the fqdn of the out-of-band interface,
	// if no hostnames/ddns is used it's fine to provide an ip-address as a string.
	// The username specified needs to have permissions to control subscriptions,
	// and the redfish protocol must be enabled in the out-of-band interface.
	client := salmongo.SalmonClient("ASDF1234.fqdn.tld", "username", "password")

	// create the subscription
	// Context can be anything you want
	//
	// Destination should point to the url of your event-listener/receiver
	// Note that subscriptions will fail unless the url includes "https"
	//
	// As of now most vendors only support the Alert and StatusChange types
	// Protocol can only be Redfish (but has to be supplied regardless)
	subscription := salmongo.Subscription{
		Context:     "foo bar",
		Destination: "https://my-listener.fqdn.tld",
		EventTypes:  []string{"Alert", "StatusChange"},
		Protocol:    "Redfish",
	}

	// Instruct the client to create the request object and send the request.
	// Handle any errors
	s, err := client.CreateSubscription(&subscription)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	out, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
