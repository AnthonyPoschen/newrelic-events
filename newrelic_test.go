package newrelic

import "testing"

var packetFormingData = []struct {
	input    []map[string]interface{}
	expected string
}{
	{
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
		expected: `{"key1":"Value1","key2",2},{"key3":"Value3","key4",4}`,
	},
}

func TestPacketForming(t *testing.T) {

	// for k, v := range packetFormingData {
	//
	// }

}
