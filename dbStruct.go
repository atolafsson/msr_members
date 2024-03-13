package main

type Member struct {
	ID       int
	Name     string
	Rank     int
	Since    string
	Phone    string
	Address  string
	City     string
	State    string
	Zip      int
	Email    string
	BirthDay string
	Notes    string
	Status   int
}

type MemberList []Member

type MemberS struct {
	Name  string
	Email string
	Phone string
}

type MemberSList []MemberS
