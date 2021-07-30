package vo

type Resort_req struct {
	MoveNodeId int64 `json:"move_node_id"`
	NewPaid int64 `json:"new_paid"`
	NewSort string `json:"new_sort"`
}
