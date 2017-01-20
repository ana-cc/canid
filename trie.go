package canid

import (
	"net"
)

// Trie for storing fast lookups of information by prefix.
// Not yet tested or integrated with canid.

type Trie struct {
	sub  [2]*Trie
	data interface{}
}

// Return the prefix and data associated with a given IP address in the trie
func (t *Trie) Find(addr net.IP) (pfx net.IPNet, data interface{}, ok bool) {

	addrmasks := [8]byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}
	netmask := make([]byte, len(addr))
	current := t

	// and iterate
	for pfxlen := 0; pfxlen < (len(addr) * 8); pfxlen++ {
		// return data if the current trie node is a leaf
		if current.data != nil {
			cnetmask := net.IPMask(netmask)
			return net.IPNet{addr.Mask(cnetmask), cnetmask}, current.data, true
		}

		// otherwise determine whether to go right or left
		var branch int
		if addr[pfxlen/8]&addrmasks[pfxlen%8] == 0 {
			branch = 0
		} else {
			branch = 1
		}

		current = current.sub[branch]

		// stop searching if nil
		if current == nil {
			break
		}

		// and move to the next bit
		netmask[pfxlen/8] |= addrmasks[pfxlen%8]
	}

	return net.IPNet{}, nil, false
}

// Add a prefix to the trie and associate some data with it

func (t *Trie) Add(pfx net.IPNet, data interface{}) {
	addrmasks := [8]byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}

	current := t
	subidx := 0

	// first search to the bottom of the trie, creating nodes as necessary
	for i := 0; pfx.Mask[i/8]&addrmasks[i%8] > 0; i++ {

		if pfx.IP[i/8]&addrmasks[i%8] == 0 {
			subidx = 0
		} else {
			subidx = 1
		}

		if current.sub[subidx] == nil {
			current.sub[subidx] = new(Trie)
		}
		current = current.sub[subidx]
	}

	/* now add data */
	current.data = data

}
