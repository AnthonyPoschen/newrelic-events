package newrelicEvents

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var packetFormingData = []struct {
	Event    string
	input    []map[string]interface{}
	expected string
}{
	{
		Event: "MyEvent",
		input: []map[string]interface{}{
			{
				"key1": "Value1",
				"key2": 2,
			},
			{
				"key3": "Value3",
				"key4": 4,
			},
		},
		expected: `[{"key1":"Value1","key2":2,"eventType":"MyEvent"},{"key3":"Value3","key4":4,"eventType":"MyEvent"}]`,
	},
}

func TestPacketForming(t *testing.T) {
	for k, v := range packetFormingData {
		nr := New("", "")
		for _, i := range v.input {
			nr.RecordEvent(v.Event, i)
		}
		var expMap []map[string]interface{}
		err := json.Unmarshal([]byte(v.expected), &expMap)
		if err != nil {
			t.Fatalf("Test: %d malformed expected", k)
		}
		var actualMap []map[string]interface{}
		err = json.Unmarshal([]byte(fmt.Sprintf("[%s]", nr.data.Data)), &actualMap)
		if err != nil {
			t.Fatalf("Test: %d Actual has invalidJson", k)
		}
		if !reflect.DeepEqual(expMap, actualMap) {
			t.Fatalf("Test: %d\nE:%s\nA:%s\n", k, v.expected, nr.data.Data)
		}
	}

}

var recordEventBadInputsData = []struct {
	Name  string
	Input map[string]interface{}
}{
	{
		Name:  "",
		Input: map[string]interface{}{"valid": "Input"},
	},
	{
		Name:  "A Name",
		Input: map[string]interface{}{"Time to break Json Marshal": make(chan int)},
	}, {
		Name:  "nil input",
		Input: nil,
	},
}

func TestRecordEventBadInputs(t *testing.T) {
	nr := New("", "")
	for k, v := range recordEventBadInputsData {
		if err := nr.RecordEvent(v.Name, v.Input); err == nil {
			t.Fatalf("Test: %d has no error", k)
		}
	}

}
