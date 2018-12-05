package clip

import (
    "testing"
)

func TestParse(t *testing.T) {
    var v8 int8
    if Parse(&v8, "127") != nil {
        t.Error(`Parse(int8, "255") error`)
    }
    if Parse(&v8, "257") == nil {
        t.Error(`Parse(int8, "257") no error`)
    }
}
