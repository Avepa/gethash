package usecase

import (
	"hash"
	"testing"

	"github.com/avepa/gethash/file"
)

func TestNewController(t *testing.T) {
	ctrl := NewController(nil)
	if ctrl == nil {
		t.Fatal("The NewController should have created an object")
	}
}

func TestnewControllerFileHash(t *testing.T) {
	type inputData struct {
		maxProc  int
		nameHash string
		ctrlRepo file.Repository
	}

	type expectedData struct {
		newHash func() hash.Hash
		err     error
	}

	type structTestData struct {
		input        inputData
		expectedData expectedData
	}

	testData := []structTestData{
		{
			input: inputData{
				maxProc:  2,
				nameHash: "",
			},
			expectedData: expectedData{
				// newHash:
			},
		},
	}

	for _, td := range testData {
		ctrl, err := newControllerFileHash(td.input.maxProc, td.input.nameHash, td.input.ctrlRepo)
		if err != td.expectedData.err {
			t.Errorf("Wrong error received.\nInput string hash: %v;\nExpected err: %v;\nReceived err: %v.",
				td.input, td.expectedData.err, err)
		} else if err == nil {
			if ctrl == nil {
				t.Fatal("The NewController should have created an object")
				continue
			}
			if cap(ctrl.nameFile) != td.input.maxProc {
				t.Errorf("The filePath channel has the wrong capacitance. Expectation: %v, received: %v.",
					td.input.maxProc, cap(ctrl.nameFile))
			}
			if cap(ctrl.hashFile) != td.input.maxProc {
				t.Errorf("The filePath channel has the wrong capacitance. Expectation: %v, received: %v",
					td.input.maxProc, cap(ctrl.hashFile))
			}
		}
	}

}
