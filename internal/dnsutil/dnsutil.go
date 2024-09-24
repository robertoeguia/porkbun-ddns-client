package dnsutil

import (
	"errors"
	"fmt"
	"log"

	"github.com/miekg/dns"
)

var (
	nameServer = "8.8.8.8"
)

func SetNameserver(ns string) {
	if ns == "" {
		log.Println("Warning: provided nameserver empty. Using default 8.8.8.8")
	}

	nameServer = ns
}

func GetARecord(host string) (string, error) {

	if host == "" {
		return "", errors.New("a valid hostname must be provided")
	}

	var message = new(dns.Msg)
	message.SetQuestion(dns.Fqdn(host), dns.TypeA)

	response := sendRequest(message)
	// (*dns.A) is type assertion. https://go.dev/ref/spec#Type_assertions
	// It asserts that in.Answer[0] is not nil and that the value stored is of type
	// dns.A
	// the ok, show below is a untyped boolean value. The value of ok is true,
	// if the assertion holds. Otherwise it is false.
	if _, ok := response.(*dns.A); !ok {
		log.Fatalf("Error record received is not correct type")
	}

	return response.(*dns.A).A.String(), nil
}

func sendRequest(msg *dns.Msg) dns.RR {
	var client = new(dns.Client)

	in, _, err := client.Exchange(msg, fmt.Sprintf("%v:53", nameServer))
	if err != nil {
		log.Fatalln(err)
	}

	return in.Answer[0]
}
