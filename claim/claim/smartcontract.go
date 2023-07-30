package claim

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	jsoniter "github.com/json-iterator/go"
)

// SmartContract ...
type SmartContract struct {
	contractapi.Contract
}

// InitLedger добавляет базовые заявки с разными статусами
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}
	claims := []Claim{
		{
			ID:        "claim1",
			Username:  "user1",
			Status:    StatusClaimClarification,
			Content:   "some content",
			UpdatedAt: now,
		},
		{
			ID:        "claim2",
			Username:  "user2",
			Status:    StatusClaimRevoke,
			Content:   "some content",
			UpdatedAt: now,
		},
		{
			ID:        "claim3",
			Username:  "user3",
			Status:    StatusClaimSatisfy,
			Content:   "some content",
			UpdatedAt: now,
		},
		{
			ID:        "claim4",
			Username:  "user4",
			Status:    StatusClaimOpen,
			Content:   "some content",
			UpdatedAt: now,
		},
		{
			ID:        "claim5",
			Username:  "user5",
			Status:    StatusClaimReject,
			Content:   "some content",
			UpdatedAt: now,
		},
		{
			ID:        "claim6",
			Username:  "user6",
			Status:    StatusClaimPending,
			Content:   "some content",
			UpdatedAt: now,
		},
	}

	for _, claim := range claims {
		assetJSON, err := jsoniter.Marshal(claim)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(claim.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ClaimExists существует ли заявка
func (s *SmartContract) ClaimExists(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	claim, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return claim != nil, nil
}

// CreateClaim Создание заявки
func (s *SmartContract) CreateClaim(ctx contractapi.TransactionContextInterface, ID string, username string, content string) error {
	exists, err := s.ClaimExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the claim %s already exists", ID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claim, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  username,
		Status:    StatusClaimOpen,
		Content:   content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claim)
}

// UpdateClaim Обновление заявки по номеру
func (s *SmartContract) UpdateClaim(ctx contractapi.TransactionContextInterface, ID string, content string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimOpen && claim.Status != StatusClaimClarification {
		return fmt.Errorf("you cannot update claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    claim.Status,
		Content:   content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}

// GetClaim Получение одной заявки по номеру
func (s *SmartContract) GetClaim(ctx contractapi.TransactionContextInterface, ID string) (*Claim, error) {
	claimBytes, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if claimBytes == nil {
		return nil, fmt.Errorf("claim %s does not exist", ID)
	}

	var claim Claim
	err = jsoniter.Unmarshal(claimBytes, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

// ListClaim Получение списка всех заявок
func (s *SmartContract) ListClaim(ctx contractapi.TransactionContextInterface) ([]*Claim, error) {
	resultsIter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIter.Close()

	var claims []*Claim
	for resultsIter.HasNext() {
		res, err := resultsIter.Next()
		if err != nil {
			return nil, err
		}

		var claim Claim
		err = jsoniter.Unmarshal(res.Value, &claim)
		if err != nil {
			return nil, err
		}
		claims = append(claims, &claim)
	}

	return claims, nil
}

// DoClaim Исполнение заявки
func (s *SmartContract) DoClaim(ctx contractapi.TransactionContextInterface, ID string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimOpen && claim.Status != StatusClaimClarification {
		return fmt.Errorf("you cannot do claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    StatusClaimPending,
		Content:   claim.Content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}

// RejectClaim Отказ в исполнении заявки
func (s *SmartContract) RejectClaim(ctx contractapi.TransactionContextInterface, ID string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimPending {
		return fmt.Errorf("you cannot reject claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    StatusClaimReject,
		Content:   claim.Content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}

// SatisfyClaim Удовлетворение заявки
func (s *SmartContract) SatisfyClaim(ctx contractapi.TransactionContextInterface, ID string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimPending {
		return fmt.Errorf("you cannot satisfy claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    StatusClaimSatisfy,
		Content:   claim.Content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}

// ClarificationClaim Уточнение данных по заявке
func (s *SmartContract) ClarificationClaim(ctx contractapi.TransactionContextInterface, ID string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimPending {
		return fmt.Errorf("you cannot clarification claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    StatusClaimClarification,
		Content:   claim.Content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}

// RevokeClaim Отзыв заявки
func (s *SmartContract) RevokeClaim(ctx contractapi.TransactionContextInterface, ID string) error {
	claim, err := s.GetClaim(ctx, ID)
	if err != nil {
		return err
	}

	if claim.Status != StatusClaimPending {
		return fmt.Errorf("you cannot revoke claim at status='%s'", claim.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	claimBytes, err := jsoniter.Marshal(Claim{
		ID:        ID,
		Username:  claim.Username,
		Status:    StatusClaimRevoke,
		Content:   claim.Content,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, claimBytes)
}
