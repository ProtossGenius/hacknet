package hacknetitf

// ServerForClientItf server for client.
type ServerForClientItf interface {
	// AcceptHacker Hacker ask for login. extraData is hacker's detail info.
	AcceptHacker(ip string, port int, email, extraData string) (result string)
	// Hack connect to another Hacker's computer.
	Hack(hackerEmail, hackerHost, targetEmail, extraData string) (result string)
}
