package claim

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

// StatusClaim Статус заявки
type StatusClaim int

var (
	// StatusClaimOpen Открыта
	StatusClaimOpen StatusClaim = 1
	// StatusClaimPending На рассмотрении
	StatusClaimPending StatusClaim = 2
	// StatusClaimReject Отказ
	StatusClaimReject StatusClaim = 3
	// StatusClaimSatisfy Удовлетворена
	StatusClaimSatisfy StatusClaim = 4
	// StatusClaimClarification Уточнение данных
	StatusClaimClarification StatusClaim = 5
	// StatusClaimRevoke Отозвана
	StatusClaimRevoke StatusClaim = 6
)

func (c StatusClaim) String() string {
	switch c {
	case StatusClaimOpen:
		return "открыта"
	case StatusClaimPending:
		return "на рассмотрении"
	case StatusClaimReject:
		return "отказ"
	case StatusClaimSatisfy:
		return "удовлетворена"
	case StatusClaimClarification:
		return "уточнение данных"
	case StatusClaimRevoke:
		return "отозвана"
	default:
		return "неизвестный статус"
	}
}

type Claim struct {
	ID        string               `json:"id"`
	Username  string               `json:"username"`
	Status    StatusClaim          `json:"status"`
	Content   string               `json:"content"`
	UpdatedAt *timestamp.Timestamp `json:"updated_at"`
}
