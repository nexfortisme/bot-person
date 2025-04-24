package enums

type RewardStatus int

const (
	TOO_EARLY RewardStatus = iota
	AVAILABLE
	MISSED
)

func (r RewardStatus) String() string {
	return [...]string{"TOO_EARLY", "AVAILABLE", "MISSED"}[r]
}
