package dynalist

import (
	"github.com/xconstruct/dynalister/pkg/api"
)

type DocumentNode api.Node

func (n DocumentNode) NodeID() string {
	return n.ID
}

func (n DocumentNode) NodeChildren() []string {
	return n.Children
}
