package webhook

type Settings struct {
	TargetURL     string
	RequestAmount int
	RPS           int
	WorkersNumber int
}
