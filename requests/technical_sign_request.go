package requests

type TechnicalSignRequest struct {
	TechnicalMsgHash string  `json:"technicalMsgHash"`
	TechnicalAddress string  `json:"technicalAddress" valid:"validHex40~ErrorAddressNotFollowingRegex"`
	TechnicalSig     SignDto `json:"technicalSig"`
}

func (tsr *TechnicalSignRequest) GetTechnicalMsgHash() string {
	return tsr.TechnicalMsgHash
}

func (tsr *TechnicalSignRequest) SetTechnicalMsgHash(TechnicalMsgHash string) {
	tsr.TechnicalMsgHash = TechnicalMsgHash
}

func (tsr *TechnicalSignRequest) GetTechnicalAddress() string {
	return tsr.TechnicalAddress
}

func (tsr *TechnicalSignRequest) SetTechnicalAddress(TechnicalAddress string) {
	tsr.TechnicalAddress = TechnicalAddress
}

func (tsr *TechnicalSignRequest) GetTechnicalSig() SignDto {
	return tsr.TechnicalSig
}

func (tsr *TechnicalSignRequest) SetTechnicalSig(TechnicalSig SignDto) {
	tsr.TechnicalSig = TechnicalSig
}
