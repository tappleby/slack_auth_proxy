package main


func main() {
	slack := NewSlackApi("xoxp-2174218295-2174484692-2250524839-0e928d")
	slack.Groups.List()
}
