package dag

import (
	"fmt"
	dagLib "github.com/heimdalr/dag"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types/meta"
)

type nodeId string

func (n nodeId) ID() string {
	return string(n)
}

func TopologySort(item *model.Dag) ([]meta.Step, error) {
	d := dagLib.NewDAG()
	nodeMap := make(map[string]*model.Node)
	for i, node := range item.Nodes {
		_, err := d.AddVertex(nodeId(node.Id))
		if err != nil {
			return nil, err
		}
		nodeMap[node.Id] = item.Nodes[i]
	}
	for _, edge := range item.Edges {
		err := d.AddEdge(edge.Source, edge.Target)
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("dag %s: %s", item.UID, d.String())

	baseRoots := d.GetRoots()
	roots := baseRoots
	have := make(map[string]struct{}, len(item.Nodes))
	var result []meta.Step
	for {
		if len(roots) == 0 {
			break
		}
		for id := range roots {
			if _, ok := have[id]; ok {
				continue
			}
			parents, err := d.GetParents(id)
			if err != nil {
				return nil, err
			}
			dependNodeId := make([]string, 0)
			for pid := range parents {
				dependNodeId = append(dependNodeId, pid)
			}
			state := model.StepCreated
			_, ok := baseRoots[id]
			if ok {
				state = model.StepReady
			}
			result = append(result, meta.Step{
				DagUID: item.UID,
				NodeId: id,
				Depend: dependNodeId,
				State:  state,
			})
			have[id] = struct{}{}
		}

		children := make(map[string]interface{})
		for id := range roots {
			items, err := d.GetChildren(id)
			if err != nil {
				return nil, err
			}
			for k, v := range items {
				children[k] = v
			}
		}
		roots = children
	}

	return result, nil
}
