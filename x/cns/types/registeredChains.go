package types

func NewRegisteredChains() *RegisteredChains {
	return &RegisteredChains{Chains: make(map[string]string, 0)}
}

func DefaultRegisteredChains() *RegisteredChains {
	return NewRegisteredChains()
}

func (rc *RegisteredChains) Validate() error {
	return nil
}
