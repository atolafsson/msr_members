package main

type Member struct {
	id       int
	name     string
	email    string
	nickName string
	prospect int
	address  string
	zip      int
	city     string
	notes    string
}

type MemberList []Member

type MemberS struct {
	name     string
	email    string
	nickName string
}

type MemberSList []MemberS
