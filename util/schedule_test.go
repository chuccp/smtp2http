package util

import "testing"

func TestParserCron(t *testing.T) {

	err := ParserCron("0 0 0/1 * * ?")
	t.Error(err)

}
