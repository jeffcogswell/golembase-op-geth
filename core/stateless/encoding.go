// Copyright 2024 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package stateless

import (
	"io"

	"github.com/jeffcogswell/golembase-op-geth/common/hexutil"
	"github.com/jeffcogswell/golembase-op-geth/core/types"
	"github.com/jeffcogswell/golembase-op-geth/crypto"
	"github.com/jeffcogswell/golembase-op-geth/rlp"
)

// toExtWitness converts our internal witness representation to the consensus one.
func (w *Witness) toExtWitness() *extWitness {
	ext := &extWitness{
		Headers: w.Headers,
	}
	ext.Codes = make([][]byte, 0, len(w.Codes))
	for code := range w.Codes {
		ext.Codes = append(ext.Codes, []byte(code))
	}
	ext.State = make([][]byte, 0, len(w.State))
	for node := range w.State {
		ext.State = append(ext.State, []byte(node))
	}
	return ext
}

// fromExtWitness converts the consensus witness format into our internal one.
func (w *Witness) fromExtWitness(ext *extWitness) error {
	w.Headers = ext.Headers

	w.Codes = make(map[string]struct{}, len(ext.Codes))
	for _, code := range ext.Codes {
		w.Codes[string(code)] = struct{}{}
	}
	w.State = make(map[string]struct{}, len(ext.State))
	for _, node := range ext.State {
		w.State[string(node)] = struct{}{}
	}
	return nil
}

// EncodeRLP serializes a witness as RLP.
func (w *Witness) EncodeRLP(wr io.Writer) error {
	return rlp.Encode(wr, w.toExtWitness())
}

// DecodeRLP decodes a witness from RLP.
func (w *Witness) DecodeRLP(s *rlp.Stream) error {
	var ext extWitness
	if err := s.Decode(&ext); err != nil {
		return err
	}
	return w.fromExtWitness(&ext)
}

// extWitness is a witness RLP encoding for transferring across clients.
type extWitness struct {
	Headers []*types.Header
	Codes   [][]byte
	State   [][]byte
}

// ExecutionWitness is a witness json encoding for transferring across clients
// in the future, we'll probably consider using the extWitness format instead for less overhead.
// currently we're using this format for compatibility with reth and also for simplicity in terms of parsing.
type ExecutionWitness struct {
	Headers []*types.Header   `json:"headers"`
	Codes   map[string]string `json:"codes"`
	State   map[string]string `json:"state"`
}

func transformMap(in map[string]struct{}) map[string]string {
	out := make(map[string]string, len(in))
	for item := range in {
		bytes := []byte(item)
		key := crypto.Keccak256Hash(bytes).Hex()
		out[key] = hexutil.Encode(bytes)
	}
	return out
}

// ToExecutionWitness converts a witness to an execution witness format that is compatible with reth.
// keccak(node) => node
// keccak(bytecodes) => bytecodes
func (w *Witness) ToExecutionWitness() *ExecutionWitness {
	return &ExecutionWitness{
		Headers: w.Headers,
		Codes:   transformMap(w.Codes),
		State:   transformMap(w.State),
	}
}
