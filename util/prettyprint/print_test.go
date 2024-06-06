package prettyprint

import "testing"

func TestPrintTable(t *testing.T) {

	header := []string{"name", "gender", "hobby"}
	data := [][]string{
		[]string{"Alice", "female", "sing"},
		[]string{"Bob", "male", "draw"},
		[]string{"Cathy", "female", "fish"},
	}

	PrintTable(header, data)
}
