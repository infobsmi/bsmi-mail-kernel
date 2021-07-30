package vo

type Nav_item struct {
	Name string `json:"name"`
	Id uint64 `json:"id"`
	Children []Nav_item `json:"children"`
	LoadOnDemand bool `json:"load_on_demand"`
}
