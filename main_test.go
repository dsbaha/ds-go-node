package main

import (
	"testing"
)

const (
	//createJobResponse = "CREATE_JOBS,dabe4605c021f09509afba3e78493fec4d05f200,6"
	lastHash = "dabe4605c021f09509afba3e78493fec4d05f200"
	//notask = "NO_TASK,,0"
)

func TestCreateJobs(t *testing.T) {
	cj := &CreateJob{
		LastHash: lastHash,
		Difficulty: uint64(6),
	}

	err := cj.createJobs()
	if err != nil {
		t.Errorf("err %s", err)
	}

	_, err = cj.marshal()
	if err != nil {
		t.Errorf("err %s", err)
	}

}

func TestMakeJob(t *testing.T) {
	j := Job{
		LastHash: lastHash,
	}

	err := makeJob(&j, uint64(1))
	if err != nil {
		t.Errorf("err %s", err)
	}

	if j.ExpectedHash == "" {
		t.Errorf("empty expected hash")
	}
}

func BenchmarkMakeJob(b *testing.B) {

	j := Job{
		LastHash: lastHash,
	}

	for i := 0; i < b.N; i++ {
		err := makeJob(&j, uint64(1))
		if err != nil {
			b.Errorf("err %s", err)
		}
	}
}
