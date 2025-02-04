package transaction

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoSaveAndLoad(t *testing.T) {
	originalMap := map[FuncSelector]string{{1, 1, 1, 1}: "selector"}
	originalMapBytes, err := json.Marshal(originalMap)
	assert.Nil(t, err)
	fmt.Printf("%v\n", originalMap)

	resultMap := map[FuncSelector]string{}
	err = json.Unmarshal(originalMapBytes, &resultMap)
	assert.Nil(t, err)
	fmt.Printf("%v\n", resultMap)
	assert.Equal(t, originalMap, resultMap)

	originalInfo := &monitored{
		Selector: &FuncSelector{2, 2, 2, 2},
	}
	fmt.Printf("%v\n", originalInfo.Selector)

	originalInfoBytes, err := json.Marshal(originalInfo)
	assert.Nil(t, err)

	resultInfo := &monitored{}
	err = json.Unmarshal(originalInfoBytes, resultInfo)
	assert.Nil(t, err)
	fmt.Printf("%v\n", resultInfo.Selector)
	assert.Equal(t, originalInfo, resultInfo)
}
