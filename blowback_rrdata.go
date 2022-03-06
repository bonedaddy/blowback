package blowback

import (
	"encoding/hex"
	"errors"

	"github.com/miekg/dns"
)

const typeDataBundle = 65281

type DataBundle struct {
	Target string
	Extra  string
}

func (db *DataBundle) String() string { return db.Target + " " + db.Extra }
func (db *DataBundle) Len() int       { return len(db.Target) + len(db.Extra)/2 }

func (db *DataBundle) Pack(buf []byte) (int, error) {
	off, err := dns.PackDomainName(db.Target, buf, 0, nil, false)
	if err != nil {
		return off, err
	}
	h, err := hex.DecodeString(db.Extra)
	if err != nil {
		return off, err
	}
	if off+hex.DecodedLen(len(db.Extra)) > len(buf) {
		return len(buf), errors.New("overflow packing hex")
	}
	copy(buf[off:off+hex.DecodedLen(len(db.Extra))], h)
	off += hex.DecodedLen(len(db.Extra))
	return off, nil
}

func (db *DataBundle) Unpack(buf []byte) (int, error) {
	s, off, err := dns.UnpackDomainName(buf, 0)
	if err != nil {
		return len(buf), err
	}
	db.Target = s
	s = hex.EncodeToString(buf[off:])
	db.Extra = s
	return len(buf), nil
}

func (db *DataBundle) Parse(sx []string) error {
	if len(sx) < 2 {
		return errors.New("need at least 2 pieces of rdata")
	}
	db.Target = sx[0]
	if _, ok := dns.IsDomainName(db.Target); !ok {
		return errors.New("bad DataBundle Target")
	}
	// Hex data can contain spaces.
	for _, s := range sx[1:] {
		db.Extra += s
	}
	return nil
}

func (db *DataBundle) Copy(dest dns.PrivateRdata) error {
	m1, ok := dest.(*DataBundle)
	if !ok {
		return dns.ErrRdata
	}
	m1.Target = db.Target
	m1.Extra = db.Extra
	return nil
}

func Register() {
	dns.PrivateHandle("DataBundle", typeDataBundle, func() dns.PrivateRdata { return new(DataBundle) })
	defer ClosePrivateRrdata()
}

func ClosePrivateRrdata() {
	dns.PrivateHandleRemove(typeDataBundle)
}
