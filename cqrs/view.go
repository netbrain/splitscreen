package cqrs

type ViewType string

type ViewRepository struct {
	views map[ViewType]interface{}
}

func NewViewRepository() *ViewRepository {
	return &ViewRepository{views: map[ViewType]interface{}{}}
}

func (v *ViewRepository) RegisterView(typ ViewType, val interface{}) {
	if _, ok := v.views[typ]; ok {
		panic("view type already registered!")
	}
	v.views[typ] = val
}

func (v *ViewRepository) GetView(typ ViewType) interface{} {
	if _, ok := v.views[typ]; !ok {
		return nil
	}
	return v.views[typ]
}
