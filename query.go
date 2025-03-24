package ecs



type Query interface {
    Apply(ecs *ECS) []EntityHandle
}

// Query for an exact Archetype
type ArchetypeQuery ComponentHash
func NewArchQuery(comps ...ComponentHandle) Query {
    query := ArchetypeQuery(CreateComponentHash(comps...))
    return &query
}
func (q *ArchetypeQuery) Apply(ecs *ECS) []EntityHandle {
    handle := ecs.archetypeQuery[ComponentHash(*q)]
    arch := ecs.archetypes[handle]
    return arch.GetEntities()
}

// Query for any entity containing at least these components
type ComponentQuery struct {
    Comps   []ComponentHandle
}

func (q *ComponentQuery) Apply(ecs *ECS) []EntityHandle {
    if len(q.Comps) == 0 { return nil }

    minSet := ecs.componentIndex[q.Comps[0]]
    for _, c := range q.Comps[1:] {
        if len(ecs.componentIndex[c]) < len(minSet) {
            minSet = ecs.componentIndex[c]
        }
    }

    var ents []EntityHandle
    for a := range minSet {
        has := true
        for _, c := range q.Comps {
            if !ecs.componentIndex[c].Has(a) {
                has = false
                break
            }
        }

        if !has { continue }
        arch := ecs.archetypes[a]
        ents = append(ents, arch.GetEntities()...)
    }

    return ents
}













