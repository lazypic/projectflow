package main

import (
	"fmt"
)

// Project 는 사용자 정보를 다루는 자료구조이다.
type Project struct {
	ID            string // 프로젝트 코드
	Budget        uint64 // 총 버짓
	MonetaryUnit  string // 화폐단위
	StartDate     string // 시작일
	EndDate       string // 마감일
	ProjectStatus string // 프로젝트 상태
	UpdateDate    string // 시작일
}

// Milestone 은 프로젝트 마일스톤 자료구조이다.
type Milestone struct {
	Title string
	Dday  string
}

func (p Project) String() string {
	return fmt.Sprintf(`
ID: %s Status: %s
Budget: %d
Start: %s
End: %s`,
		p.ID,
		p.ProjectStatus,
		p.Budget,
		p.StartDate,
		p.EndDate,
	)
}
