package csv_test

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []testMatrixItem{{
		name: "1",
		args: args{s: "foo,bar"},
		want: [][]string{{
			"foo", "bar",
		}},
		wantErr: false,
	}, {
		name: "2",
		args: args{s: `"foo","bar"`},
		want: [][]string{{
			"foo", "bar",
		}},
		wantErr: false,
	}, {
		name: "3",
		args: args{s: "foo,bar\r\n1,2"},
		want: [][]string{{
			"foo", "bar",
		}, {
			"1", "2",
		}},
		wantErr: false,
	}, {
		name: "4",
		args: args{s: `"foo","bar"` + "\r\n" + `"1","2"`},
		want: [][]string{{
			"foo", "bar",
		}, {
			"1", "2",
		}},
		wantErr: false,
	}, {
		name:    "5",
		args:    args{s: `"foo","bar"` + "\r\n" + `"1","2`},
		want:    nil,
		wantErr: true,
	}}

	runMatrixParse(t, tests)
}
