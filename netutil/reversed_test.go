package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ipv4RevGood   = `1.0.0.127.in-addr.arpa`
	ipv4RevGoodUp = `1.0.0.127.In-Addr.Arpa`

	ipv4RevGoodUnspecified = `0.0.0.0.in-addr.arpa`

	ipv4Missing = `.0.0.127.in-addr.arpa`
	ipv4Char    = `1.0.z.127.in-addr.arpa`
)

const (
	ipv6RevZeroes = `0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0`
	ipv6Suffix    = ipv6RevZeroes + `.ip6.arpa`

	ipv6RevGood   = `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevGoodUp = `4.3.2.1.D.C.B.A.0.0.0.0.0.0.0.0.` + ipv6Suffix

	ipv6RevGoodUnspecified = `0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.` + ipv6Suffix

	ipv6RevCharHi  = `4.3.2.1.d.c.b.a.0.z.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevCharLo  = `4.3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevDots    = `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6RevZeroes + `..ip6.arpa`
	ipv6RevLen     = `3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevMany    = `4.3.2.1.dbc.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevMissing = `.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
	ipv6RevSpace   = `4.3.2.1.d.c.b.a. .0.0.0.0.0.0.0.` + ipv6Suffix
)

func TestUnreverseAddr(t *testing.T) {
	ip4 := net.IP{127, 0, 0, 1}
	ip6 := net.IP{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xab, 0xcd, 0x12, 0x34,
	}

	testCases := []struct {
		name       string
		in         string
		wantErrMsg string
		wantErrAs  interface{}
		want       net.IP
	}{{
		name:       "good_ipv4",
		in:         ipv4RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip4,
	}, {
		name:       "good_ipv4_fqdn",
		in:         ipv4RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip4,
	}, {
		name:       "good_ipv4_case",
		in:         ipv4RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip4,
	}, {
		name: "bad_ipv4_missing",
		in:   ipv4Missing,
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": label is empty`,
		wantErrAs: new(errors.Error),
		want:      nil,
	}, {
		name: "bad_ipv4_char",
		in:   ipv4Char,
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`bad ipv4 address "1.0.z.127"`,
		wantErrAs: new(*netutil.AddrError),
		want:      nil,
	}, {
		name:       "good_ipv6",
		in:         ipv6RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip6,
	}, {
		name:       "good_ipv6_fqdn",
		in:         ipv6RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip6,
	}, {
		name:       "good_ipv6_case",
		in:         ipv6RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       ip6,
	}, {
		name: "bad_ipv6_many",
		in:   ipv6RevMany,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMany + `": ` +
			`not a full reversed ip address`,
		wantErrAs: new(*netutil.AddrError),
		want:      nil,
	}, {
		name: "bad_ipv6_missing",
		in:   ipv6RevMissing,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMissing + `": ` +
			`bad domain name label "": label is empty`,
		wantErrAs: new(errors.Error),
		want:      nil,
	}, {
		name: "bad_ipv6_char_lo",
		in:   ipv6RevCharLo,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_char_hi",
		in:   ipv6RevCharHi,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_dots",
		in:   ipv6RevDots,
		wantErrMsg: `bad arpa domain name "` + ipv6RevDots + `": ` +
			`bad domain name label "": label is empty`,
		wantErrAs: new(errors.Error),
		want:      nil,
	}, {
		name: "bad_ipv6_len",
		in:   ipv6RevLen,
		wantErrMsg: `bad arpa domain name "` + ipv6RevLen + `": ` +
			`bad arpa domain name length 70, allowed: 72`,
		wantErrAs: new(*netutil.LengthError),
		want:      nil,
	}, {
		name: "bad_ipv6_space",
		in:   ipv6RevSpace,
		wantErrMsg: `bad arpa domain name "` + ipv6RevSpace + `": ` +
			`bad domain name label " ": bad domain name label rune ' '`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name:       "not_a_reversed_ip",
		in:         "1.2.3.4",
		wantErrMsg: `bad arpa domain name "1.2.3.4": not a full reversed ip address`,
		wantErrAs:  new(errors.Error),
		want:       nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := netutil.IPFromReversedAddr(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.To16(), ip.To16())
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestIPToReversedAddr(t *testing.T) {
	ip4 := net.IP{127, 0, 0, 1}
	ip6 := net.IP{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xab, 0xcd, 0x12, 0x34,
	}

	testCases := []struct {
		name       string
		want       string
		wantErrMsg string
		wantErrAs  interface{}
		in         net.IP
	}{{
		name:       "good_ipv4",
		want:       ipv4RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         ip4,
	}, {
		name:       "good_ipv6",
		want:       ipv6RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         ip6,
	}, {
		name:       "nil_ip",
		want:       "",
		wantErrMsg: `bad ip address "<nil>"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         nil,
	}, {
		name:       "empty_ip",
		want:       "",
		wantErrMsg: `bad ip address "<nil>"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         net.IP{},
	}, {
		name:       "unspecified_ipv4",
		want:       ipv4RevGoodUnspecified,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.IPv4zero,
	}, {
		name:       "unspecified_ipv6",
		want:       ipv6RevGoodUnspecified,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.IPv6unspecified,
	}, {
		name:       "wrong_length_ip",
		want:       "",
		wantErrMsg: `bad ip address "?0102030405"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         net.IP{1, 2, 3, 4, 5},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arpa, err := netutil.IPToReversedAddr(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, arpa)
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}
