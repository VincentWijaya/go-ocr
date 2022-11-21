package validate

import (
	"context"

	"github.com/vincentwijaya/go-ocr/internal/app/repo/member"
)

type ValidateUC interface {
}

type validateUC struct {
	memberRepo member.MemberRepo
}

func New(memberRepo member.MemberRepo) *validateUC {
	return &validateUC{
		memberRepo: memberRepo,
	}
}

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, request ValidatePlateAndOwnerRequest) error {
	return nil
}
