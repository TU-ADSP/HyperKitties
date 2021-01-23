package main

import "encoding/base64"
import "strings"

func getClientID(clientID string) string {
	decoded, _ := base64.StdEncoding.DecodeString(clientID)
	decodedString := string(decoded)[9:]
	decodedSplit := strings.Split(decodedString, ",")
	return decodedSplit[0]
}
