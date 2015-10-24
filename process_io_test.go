package procmon

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSimpleProcess(t *testing.T) {
	point, err := parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735 34816 2679 4218880 655 3141 0 0 0 0 0 1 20 0 1 0 182865 12144640 534 18446744073709551615 4194304 4729572 140730058798800 140730058796792 139881841869352 0 0 2637828 2 0 0 0 17 0 0 0 0 0 0 6826728 6830659 16642048 140730058800890 140730058800894 140730058800894 140730058801136 0`))
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0), point.user)
		assert.Equal(t, uint64(0), point.system)
	}
	point, err = parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735 34816 2679 4218880 655 3141 0 0 152 189 162 199 20 0 1 0 182865 12144640 534 18446744073709551615 4194304 4729572 140730058798800 140730058796792 139881841869352 0 0 2637828 2 0 0 0 17 0 0 0 0 0 0 6826728 6830659 16642048 140730058800890 140730058800894 140730058800894 140730058801136 0`))
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(152), point.user)
		assert.Equal(t, uint64(189), point.system)
	}
}

func TestBrokenProcess(t *testing.T) {
	_, err := parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735`))
	assert.Error(t, err)
	_, err = parseProcStat(strings.NewReader(``))
	assert.Error(t, err)
	_, err = parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735 34816 2679 4218880 655 3141 0 0 frob 0 0 1 20 0 1 0 182865 12144640 534 18446744073709551615 4194304 4729572 140730058798800 140730058796792 139881841869352 0 0 2637828 2 0 0 0 17 0 0 0 0 0 0 6826728 6830659 16642048 140730058800890 140730058800894 140730058800894 140730058801136 0`))
	assert.Error(t, err)
	_, err = parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735 34816 2679 4218880 655 3141 0 0 0 botz 0 0 20 0 1 0 182865 12144640 534 18446744073709551615 4194304 4729572 140730058798800 140730058796792 139881841869352 0 0 2637828 2 0 0 0 17 0 0 0 0 0 0 6826728 6830659 16642048 140730058800890 140730058800894 140730058800894 140730058801136 0`))
	assert.Error(t, err)
	_, err = parseProcStat(strings.NewReader(`1735 (sh) S 1734 1735 1735 34816 2679 4218880 655 3141 0 0 -2 4 0 1 20 0 1 0 182865 12144640 534 18446744073709551615 4194304 4729572 140730058798800 140730058796792 139881841869352 0 0 2637828 2 0 0 0 17 0 0 0 0 0 0 6826728 6830659 16642048 140730058800890 140730058800894 140730058800894 140730058801136 0`))
	assert.Error(t, err)
}

func TestSimpleMemory(t *testing.T) {
	point, err := parseMemStat(strings.NewReader(`2965 534 485 131 0 129 0`))
	if assert.NoError(t, err) {
		assert.Equal(t, point, uint64(534))
	}
}

func TestBrokenMemory(t *testing.T) {
	_, err := parseMemStat(strings.NewReader(`1735`))
	assert.Error(t, err)
	_, err = parseMemStat(strings.NewReader(``))
	assert.Error(t, err)
	_, err = parseMemStat(strings.NewReader(`1735 (sh) S`))
	assert.Error(t, err)
	_, err = parseMemStat(strings.NewReader(`1735 -24`))
	assert.Error(t, err)
}

func TestSimpleEverything(t *testing.T) {
	input := `cpu  2019 0 929 687424 84 1 34 0 0 0
cpu0 2019 0 929 687424 84 1 34 0 0 0
intr 166310 14 10 0 0 0 0 0 0 1 0 0 0 156 0 0 0 7407 0 0 8335 10262 6589 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 296117
btime 1445583144
processes 2833
procs_running 2
procs_blocked 0
softirq 134624 0 68894 2777 14447 5489 0 18 0 29 42970`
	point, err := parseGlobalStat(strings.NewReader(input))
	if assert.NoError(t, err) {
		assert.Equal(t, point.user, uint64(2019))
		assert.Equal(t, point.system, uint64(929))
	}
	input = `cpu0 2019 0 929 687424 84 1 34 0 0 0
intr 166310 14 10 0 0 0 0 0 0 1 0 0 0 156 0 0 0 7407 0 0 8335 10262 6589 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 296117
btime 1445583144
processes 2833
cpu  2019 0 929 687424 84 1 34 0 0 0
procs_running 2
procs_blocked 0
softirq 134624 0 68894 2777 14447 5489 0 18 0 29 42970`
	point, err = parseGlobalStat(strings.NewReader(input))
	if assert.NoError(t, err) {
		assert.Equal(t, point.user, uint64(2019))
		assert.Equal(t, point.system, uint64(929))
	}
	input = `cpu        2019 0 929 687424 84 1 34 0 0 0`
	point, err = parseGlobalStat(strings.NewReader(input))
	if assert.NoError(t, err) {
		assert.Equal(t, point.user, uint64(2019))
		assert.Equal(t, point.system, uint64(929))
	}
}

func TestBrokenEverything(t *testing.T) {
	_, err := parseGlobalStat(strings.NewReader(`cpu  2 3`))
	assert.Error(t, err)
	_, err = parseGlobalStat(strings.NewReader(``))
	assert.Error(t, err)
	_, err = parseGlobalStat(strings.NewReader(`cpu0 2019 0 929 687424 84 1 34 0 0 0
intr 166310 14 10 0 0 0 0 0 0 1 0 0 0 156 0 0 0 7407 0 0 8335 10262 6589 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 296117
btime 1445583144
processes 2833
procs_running 2
procs_blocked 0
softirq 134624 0 68894 2777 14447 5489 0 18 0 29 42970`))
	assert.Error(t, err)
	_, err = parseGlobalStat(strings.NewReader(`cpu      foo 4 21`))
	assert.Error(t, err)
}
