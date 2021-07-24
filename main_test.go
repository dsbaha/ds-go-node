package main

import (
	"testing"
)

const (
	createJobResponse = "CREATE_JOBS,dabe4605c021f09509afba3e78493fec4d05f200,6"
	lastHash = "dabe4605c021f09509afba3e78493fec4d05f200"
	notask = "NO_TASK,,0"
	testName = "testUser"
)

func TestParseJobs(t *testing.T) {
	j := &CreateJob{
		User: testName,
	}

	test := string(createJobResponse)

	err := j.parseJobs(&test)
	if err != nil {
		t.Errorf("err %s", err)
		return
	}

	if (j.User != testName) {
		t.Errorf("invalid user got %s expecting %s", j.User, testName)
		return
	}

	if (j.LastHash != lastHash) {
		t.Errorf("invalid lasthash got %s expecting %s", j.LastHash, lastHash)
		return
	}

	if (j.Difficulty != uint64(6)){
		t.Errorf("invalid difficulty got %v expected %v", j.Difficulty, uint64(6))
		return
	}
}

func TestParseJobNoTask(t *testing.T) {
	j := &CreateJob{
		User: testName,
	}

	tr := 1
	wait = &tr
	qu := true
	quiet = &qu

	test := string(notask)
	err := j.parseJobs(&test)
	if err == nil {
		t.Error("should have returned an error")
		return
	}
}

func TestCreateJobs(t *testing.T) {
	cj := &CreateJob{
		User: testName,
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

	if cj.User == "" {
		t.Error("user not set")
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
			return
		}
	}
}

func BenchmarkParseJob(b *testing.B) {
	j := &CreateJob{
		User: testName,
	}

	test := string(createJobResponse)

	for i := 0; i < b.N; i++ {
		err := j.parseJobs(&test)
		if err != nil {
			b.Errorf("err %s", err)
			return
		}
	}
}
