package types

// nolint:stylecheck,st1003
// revive:disable-next-line var-naming

func (k *Keyring) IsParty(address string) bool {
	for _, party := range k.Parties {
		if party == address {
			return true
		}
	}
	return false
}

func (k *Keyring) IsAdmin(address string) bool {
	for _, admin := range k.Admins {
		if admin == address {
			return true
		}
	}
	return false
}

func (k *Keyring) RemoveAdmin(address string) {
	for i, admin := range k.Admins {
		if admin == address {
			k.Admins = append(k.Admins[:i], k.Admins[i+1:]...)
			return
		}
	}
}

func (k *Keyring) AddAdmin(address string) {
	k.Admins = append(k.Admins, address)
}

func (k *Keyring) AddParty(address string) {
	k.Parties = append(k.Parties, address)
}

func (k *Keyring) SetKeyReqFee(fee uint64) {
	k.KeyReqFee = fee
}

func (k *Keyring) SetSigReqFee(fee uint64) {
	k.SigReqFee = fee
}

func (k *Keyring) SetStatus(status bool) {
	k.IsActive = status
}

func (k *Keyring) SetDescription(description string) {
	k.Description = description
}

func (k *Keyring) RemoveParty(address string) {
	for i, party := range k.Parties {
		if party == address {
			k.Parties = append(k.Parties[:i], k.Parties[i+1:]...)
			return
		}
	}
}

func (k *Keyring) SetPartyThreshold(threshold uint32) {
	k.PartyThreshold = threshold
}

func (k *Keyring) SetMpcDefaultBtl(btl uint64) {
	k.MpcDefaultBtl = btl
}

func (k *Keyring) SetMpcMinimumBtl(btl uint64) {
	k.MpcMinimumBtl = btl
}
