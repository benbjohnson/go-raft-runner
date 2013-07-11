package runner

import (
	"github.com/benbjohnson/go-raft"
)

type JoinCommand struct {
	Name string `json:"name"`
}

func (c *JoinCommand) CommandName() string {
	return "join"
}

func (c *JoinCommand) Apply(server *raft.Server) (interface{}, error) {
	return nil, server.AddPeer(c.Name)
}
