package html

import "testing"

func Test_currency(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "zero number", args: args{amount: 0}, want: "$0.00"},
		{name: "positive number", args: args{amount: 15010}, want: "$150.10"},
		{name: "negative number", args: args{amount: -15002}, want: "-$150.02"},
		{name: "thousands number", args: args{amount: -100055}, want: "-$1,000.55"},
		{name: "large number", args: args{amount: 123456789}, want: "$1,234,567.89"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := currency(tt.args.amount); got != tt.want {
				t.Errorf("currency() = %v, want %v", got, tt.want)
			}
		})
	}
}
