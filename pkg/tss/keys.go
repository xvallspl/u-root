// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

func defaultSymScheme() *tpm2.SymScheme {
	return &tpm2.SymScheme{
		Alg:     tpm2.AlgAES,
		KeyBits: 128,
		Mode:    tpm2.AlgCFB,
	}
}

func defaultECCParams() *tpm2.ECCParams {
	return &tpm2.ECCParams{
		Symmetric: defaultSymScheme(),
		CurveID:   tpm2.CurveNISTP256,
		Point: tpm2.ECPoint{
			XRaw: make([]byte, 32),
			YRaw: make([]byte, 32),
		},
	}
}

func loadSRK20(rwc io.ReadWriteCloser, srkPW string) (*tpm2tools.Key, error) {
	var srkAuth tpmutil.U16Bytes
	var hash [32]byte

	if srkPW != "" {
		hash = sha256.Sum256([]byte(srkPW))
	}

	srkPWBytes := bytes.NewBuffer(hash[:])
	err := srkAuth.TPMMarshal(srkPWBytes)
	if err != nil {
		return nil, err
	}
	srkTemplate := tpm2.Public{
		Type:          tpm2.AlgECC,
		NameAlg:       tpm2.AlgSHA256,
		Attributes:    tpm2.FlagStorageDefault,
		ECCParameters: defaultECCParams(),
		AuthPolicy:    srkAuth,
	}

	key, err := tpm2tools.NewCachedKey(rwc, tpm2.HandleOwner, srkTemplate, tpm2tools.SRKECCReservedHandle)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func seal12(rwc io.ReadWriteCloser, srkPW string, pcrs []int, data []byte) ([]byte, error) {
	var srkAuth [20]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}
	sealed, err := tpm.Seal(rwc, tpm.Locality(1), pcrs, data, srkAuth[:])
	if err != nil {
		return nil, fmt.Errorf("couldn't seal the data", err)
	}

	return sealed, nil
}

func reseal12(rwc io.ReadWriteCloser, srkPW string, pcrs map[int][]byte, data []byte) ([]byte, error) {
	var srkAuth [20]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}
	sealed, err := tpm.Reseal(rwc, tpm.Locality(1), pcrs, data, srkAuth[:])
	if err != nil {
		return nil, fmt.Errorf("couldn't reseal the data", err)
	}

	return sealed, nil
}

func unseal12(rwc io.ReadWriteCloser, srkPW string, sealed []byte) ([]byte, error) {
	var srkAuth [20]byte

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}
	unsealed, err := tpm.Unseal(rwc, sealed, srkAuth[:])
	if err != nil {
		return nil, fmt.Errorf("couldn't seal the data", err)
	}

	return unsealed, nil
}

func seal20(rwc io.ReadWriteCloser, srkPW string, pcrs []int, data []byte) (*proto.SealedBytes, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	sOpt := tpm2tools.SealCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}
	sealed, err := key.Seal(data, sOpt)
	if err != nil {
		return nil, err
	}

	return sealed, nil
}

func unseal20(rwc io.ReadWriteCloser, srkPW string, pcrs []int, sealed *proto.SealedBytes) ([]byte, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	cOpt := tpm2tools.CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}
	unsealed, err := key.Unseal(sealed, cOpt)
	if err != nil {
		return nil, err
	}

	return unsealed, nil
}

func reseal20(rwc io.ReadWriteCloser, srkPW string, pcrs []int, sealed *proto.SealedBytes) (*proto.SealedBytes, error) {
	key, err := loadSRK20(rwc, srkPW)
	if err != nil {
		return nil, err
	}
	cOpt := tpm2tools.CertifyCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}
	sOpt := tpm2tools.SealCurrent{
		PCRSelection: tpm2.PCRSelection{
			Hash: tpm2.AlgSHA256,
			PCRs: pcrs,
		},
	}

	sealed, err = key.Reseal(sealed, cOpt, sOpt)
	if err != nil {
		return nil, err
	}

	return sealed, nil
}
