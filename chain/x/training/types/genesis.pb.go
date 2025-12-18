package types

type GenesisState struct {
	Jobs  []Job  `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs"`
	Tasks []Task `protobuf:"bytes,2,rep,name=tasks,proto3" json:"tasks"`
}

