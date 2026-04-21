package steamid

import "testing"

func TestResolve(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    uint32
		wantErr bool
	}{
		{name: "steam64", input: "76561199090171547", want: 1129905819},
		{name: "steam64 second account", input: "76561199143563790", want: 1183298062},
		{name: "steam2 odd", input: "STEAM_0:1:564952909", want: 1129905819},
		{name: "steam2 even", input: "STEAM_0:0:591649031", want: 1183298062},
		{name: "steam3", input: "[U:1:1129905819]", want: 1129905819},
		{name: "steam3 second account", input: "[U:1:1183298062]", want: 1183298062},
		{name: "steam32 direct", input: "1129905819", want: 1129905819},
		{name: "steam32 zero", input: "0", want: 0},
		{name: "whitespace trimmed", input: "  1129905819  ", want: 1129905819},
		{name: "empty string", input: "", wantErr: true},
		{name: "garbage", input: "notanid", wantErr: true},
		{name: "steam64 too small", input: "76561197960265727", wantErr: true},
		{name: "negative sign", input: "-1", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Resolve(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("Resolve(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestToSteam64(t *testing.T) {
	cases := []struct {
		accountID uint32
		want      uint64
	}{
		{1129905819, 76561199090171547},
		{1183298062, 76561199143563790},
		{0, 76561197960265728},
	}

	for _, tc := range cases {
		got := ToSteam64(tc.accountID)
		if got != tc.want {
			t.Errorf("ToSteam64(%d) = %d, want %d", tc.accountID, got, tc.want)
		}
	}
}
