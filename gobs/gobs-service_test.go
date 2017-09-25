

package gooption

import (
	"testing"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gooption/pb"
)


func Test_Price(t *testing.T) {
	if file, err := os.Open("./testdata/PriceRequest.json"); err == nil {
		defer file.Close()
		request := &pb.PriceRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			response := NewService().Price(request)
			if response.Error != "" {
				t.Error(err)
			}

			t.Log(response)
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

func Test_Greek(t *testing.T) {
	if file, err := os.Open("./testdata/GreekRequest.json"); err == nil {
		defer file.Close()
		request := &pb.GreekRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			response := NewService().Greek(request)
			if response.Error != "" {
				t.Error(err)
			}

			t.Log(response)
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

func Test_ImpliedVol(t *testing.T) {
	if file, err := os.Open("./testdata/ImpliedVolRequest.json"); err == nil {
		defer file.Close()
		request := &pb.ImpliedVolRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			response := NewService().ImpliedVol(request)
			if response.Error != "" {
				t.Error(err)
			}

			t.Log(response)
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

