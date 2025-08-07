package action

type ActionService interface {
	save(action Action)
} 

type InmemoryActionService struct {
	actions []Action
}

func (s *InmemoryActionService) Save(action Action) {
	s.actions = append(s.actions, action)
}

func (s *InmemoryActionService) GetActions() []Action {
	copiedActions := make([]Action, len(s.actions))
	copy(copiedActions, s.actions)
	return copiedActions
}