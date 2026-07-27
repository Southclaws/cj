package actions

import "sort"

type Registry struct {
	actions map[string]Action
}

func NewRegistry() *Registry {
	return &Registry{actions: make(map[string]Action)}
}

func (r *Registry) Register(a Action) {
	r.actions[a.Name()] = a
}

func (r *Registry) Get(name string) (Action, bool) {
	a, ok := r.actions[name]
	return a, ok
}

func (r *Registry) List() []Action {
	out := make([]Action, 0, len(r.actions))
	for _, a := range r.actions {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}
