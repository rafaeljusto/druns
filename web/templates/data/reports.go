package data

import (
	"github.com/rafaeljusto/druns/core/reports"
	"github.com/rafaeljusto/druns/core/types"
)

type Reports struct {
	Logged
	Incomings []reports.Incoming
}

func NewReports(username types.Name, menu Menu, incomings []reports.Incoming) Reports {
	return Reports{
		Logged:    NewLogged(username, menu),
		Incomings: incomings,
	}
}
