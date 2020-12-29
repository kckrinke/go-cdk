package cdk

type CTypeItemList []TypeItem

func (t CTypeItemList) Index(item TypeItem) int {
	for id, ri := range t {
		if ri != nil && ri.ObjectID() == item.ObjectID() {
			return id
		}
	}
	return -1
}
