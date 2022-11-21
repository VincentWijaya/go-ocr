package validate

import "github.com/vincentwijaya/go-ocr/internal/app/repo/member"

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

func (uc *validateUC) Validate() ValidateUC {
	if !uc.memberRepo {
		panic("memberRepo is nil")
	}

	return uc
}

func (uc *validateUC) ValidatePlateAndOwner()
