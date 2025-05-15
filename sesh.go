package sesh

type Session struct {
	Name    string
	Windows []window
}

type window struct {
	Title   string
	Command string
}
