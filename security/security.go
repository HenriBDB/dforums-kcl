package security

import "crypto/sha256"

type HashSignature [28]byte

type SecurityObject struct {
	Fingerprint HashSignature
	ProofOfWork string
}

//TODO Log failed checks
func (so *SecurityObject) Verify(dataByte []byte) bool {
	expectedFingerprint := createFingerprint(dataByte, so.ProofOfWork)
	if expectedFingerprint != so.Fingerprint {
		return false
	}
	ok, err := verifyProofOfWork(so.ProofOfWork, dataByte)
	if err != nil || !ok {
		return false
	}
	return true
}

func GenSecurityObject(dataBytes []byte) (SecurityObject, error) {
	so := SecurityObject{}

	// 1 Create Proof of Work
	pow, err := newProofOfWork(dataBytes)
	if err != nil {
		return SecurityObject{}, err
	}
	so.ProofOfWork = pow

	// 2 Create the fingerprint
	so.genFingerprint(dataBytes)
	return so, nil
}

func (so *SecurityObject) genFingerprint(dataBytes []byte) {
	so.Fingerprint = createFingerprint(dataBytes, so.ProofOfWork)
}

func createFingerprint(dataBytes []byte, pow string) HashSignature {
	return sha256.Sum224(append(dataBytes, []byte(pow)...))
}
