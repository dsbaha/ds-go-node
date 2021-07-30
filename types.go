package main

type CreateJob struct {
	LastHash   string `json:"-"`
	Difficulty uint64 `json:"-"`
	User       string `json:"user"`
	Jobs       []Job  `json:"jobs"`
}

type Job struct {
	LastHash     string `json:"last_hash"`
	ExpectedHash string `json:"expected_hash"`
	Nonce        uint64 `json:"numeric_result"`
}
