package main

import "log"


func main() {
	slack := NewSlackApi("xoxp-2176048118-2176048120-2250552618-570941")

	userAuth, _ := slack.Auth.Test()

	log.Println(userAuth)

//	groups, _ := slack.Groups.List()
//	g := groups.FindName("devs")

//	if g != nil {
//		log.Printf("Found dev group id: %s", g.Id)
//
//	} else {
//		log.Println("Error finding devs group")
//	}
}
