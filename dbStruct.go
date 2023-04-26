package main

type Member struct {
	id       int
	name     string
	email    string
	nickName string
	prospect int
	address  string
	zip      string
	city     string
	notes    string
}

type MemberList []Member
