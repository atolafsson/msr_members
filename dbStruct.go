package main

type Member struct {
	ID       int
	Name     string
	Email    string
	NickName string
	Prospect int
	Address  string
	Zip      int
	City     string
	Notes    string
}

type MemberList []Member

type MemberS struct {
	name     string
	email    string
	nickName string
}

type MemberSList []MemberS
