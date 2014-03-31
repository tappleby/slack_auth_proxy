package main

import "log"

type GroupService struct {
	api *SlackApi
}

func (g *GroupService) List() {

	log.Println("Doing listing")

}
