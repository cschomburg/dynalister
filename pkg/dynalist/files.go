package dynalist

import (
	"github.com/xconstruct/dynalister/pkg/api"
)

type FileNode api.File

func (n FileNode) NodeID() string {
	return n.ID
}

func (n FileNode) NodeChildren() []string {
	return n.Children
}
