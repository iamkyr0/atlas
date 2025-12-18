package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Jobs:  []Job{},
		Tasks: []Task{},
	}
}

func (gs GenesisState) Validate() error {
	return nil
}

