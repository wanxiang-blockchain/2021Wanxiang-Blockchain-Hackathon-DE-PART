package utils

type Record struct {
	ID      string `json:"id"`
	DataUri string `json:"dataUri"`
	Capsule string `json:"capsule"`
}

type ReKey struct {
	ID string `json:"id"`
	Rk string `json:"rk"`
}

type CapsuleResult struct {
	NewCapsule string `json:"newCapsule"`
	DataUri    string `json:"dataUri"`
}
