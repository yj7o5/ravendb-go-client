package operations

import (
	"encoding/json"
	"github.com/ravendb/ravendb-go-client"
	"net/http"
)

type RemoveClusterNode struct {
	Node string `json:"Node"`
	Tag  string `json:"Tag"`
}

func NewRemovePromoteClusterNode(node string) *OperationPromoteClusterNode {
	return &OperationPromoteClusterNode{
		Node: node,
	}
}

func (operation *RemoveClusterNode) GetCommand(conventions *ravendb.DocumentConventions) (ravendb.RavenCommand, error) {
	return &removeNodeCommand{
		parent: operation,
	}, nil
}

type removeNodeCommand struct {
	ravendb.RavenCommandBase
	parent *RemoveClusterNode
}

func (c *removeNodeCommand) CreateRequest(node *ravendb.ServerNode) (*http.Request, error) {
	url := node.URL + "/admin/cluster/node?nodeTag=" + c.parent.Node
	return http.NewRequest(http.MethodDelete, url, nil)
}
func (c *removeNodeCommand) SetResponse(response []byte, fromCache bool) error {
	return json.Unmarshal(response, c.parent)
}
