package cqrs

type ManagerType string

type ManagerRepository struct {
	managers map[ManagerType]interface{}
}

func NewManagerRepository() *ManagerRepository {
	return &ManagerRepository{managers: map[ManagerType]interface{}{}}
}

func (v *ManagerRepository) RegisterManager(typ ManagerType, val interface{}) {
	if _, ok := v.managers[typ]; ok {
		panic("manager type already registered!")
	}
	v.managers[typ] = val
}

func (v *ManagerRepository) GetManager(typ ManagerType) interface{} {
	if _, ok := v.managers[typ]; !ok {
		return nil
	}
	return v.managers[typ]
}
