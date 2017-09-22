package gooption

import (
	"math"
	"testing"
)

var (
	params = map[string]float64{
		"S":       100,
		"K":       100,
		"T":       1,
		"Sigma":   0.1,
		"R":       0.01,
		"Call":    1,
		"Put":     -1,
		"atmCall": 4.485236409022086,
		"atmPut":  3.4902197839388975,
	}
)

func BenchmarkFairValue(b *testing.B) {
	for index := 0; index < b.N; index++ {
		_, err := FairValue(
			params["S"],
			params["K"],
			params["T"],
			params["Sigma"],
			params["R"],
			params["Call"])
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}
func TestPutCallParity(t *testing.T) {
	call, _ := FairValue(
		params["S"],
		params["K"],
		params["T"],
		params["Sigma"],
		params["R"],
		params["Call"])
	t.Logf("Call: %v", call)

	put, _ := FairValue(
		params["S"],
		params["K"],
		params["T"],
		params["Sigma"],
		params["R"],
		params["Put"])
	t.Logf("Put: %v", put)

	forward := params["S"] - math.Exp(-params["R"]*params["T"])*params["K"]
	t.Logf("Forward: %v", forward)

	if math.Abs(call-put-forward) > 1E-10 {
		t.Errorf("Put Call parity broken")
		t.Logf("Call - Put: %v", call-put)
	}
}

func TestWrongPutCall(t *testing.T) {
	_, err := FairValue(
		params["S"],
		params["K"],
		params["T"],
		params["Sigma"],
		params["R"],
		0)
	if err == nil {
		t.Errorf("Invalid PutCall value equal to 0")
	}

	_, err = FairValue(
		params["S"],
		params["K"],
		params["T"],
		params["Sigma"],
		params["R"],
		-2)
	if err == nil {
		t.Errorf("Invalid PutCall value equal to -2")
	}

	_, err = FairValue(
		params["S"],
		params["K"],
		params["T"],
		params["Sigma"],
		params["R"],
		2)
	if err == nil {
		t.Errorf("Invalid PutCall value equal to 2")
	}
}

func TestPutCallImpliedVol(t *testing.T) {
	ivCall, iter, err := ImpliedVol(
		params["atmCall"],
		params["S"],
		params["K"],
		params["T"],
		params["R"],
		params["Call"])
	t.Logf("atmCall iv: %v", ivCall)
	t.Logf("atmCall iv iteration: %v", iter)
	if err != nil {
		t.Errorf(err.Error())
	}
	if math.Abs(ivCall-params["Sigma"]) > 1E-10 {
		t.Errorf("atmCall iv %v should be equal to %v", ivCall, params["Sigma"])
	}

	ivPut, iter, err := ImpliedVol(
		params["atmPut"],
		params["S"],
		params["K"],
		params["T"],
		params["R"],
		params["Put"])
	t.Logf("atmPut iv: %v", ivPut)
	t.Logf("atmPut iv iteration: %v", iter)
	if err != nil {
		t.Errorf(err.Error())
	}
	if math.Abs(ivPut-params["Sigma"]) > 1E-10 {
		t.Errorf("atmPut iv %v should be equal to %v", ivPut, params["Sigma"])
	}

	if math.Abs(ivCall-ivPut) > 1E-10 {
		t.Errorf("atm vol call different from atm put")
	}
}