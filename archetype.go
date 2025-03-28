package ecs

import "slices"




type Record struct {
    Arch    Archetype
    // Different handle, Record.Id is the handle local to the Archetype
    Id      EntityHandle
}

// Archetypes can be looked at as a graph.
// Each Archetype should have some sort of mapping from ComponentHandle -> ArchetypeEdge.
// An ArchetypeEdge has 2 pointers to Archetypes, one with the Component, and one without the Component,
// These can also be nil if there is no Archetype that meets those requirements
//type ArchetypeEdge struct {
//    Add     *Archetype
//    Remove  *Archetype
//}

type ArchetypeHandle uint64
const InvalidArchetypeHandle ArchetypeHandle = ArchetypeHandle(InvalidHandle)
// If you create your own Archetype struct you should register it before creating any entities.
// You should consider creating a custom archetype if you have a common entity that will
// always have the same data layout, like a player. Doing this allows you to optimise how entity data
// is stored and accessed. If you do not do this the ECS will default to a GenericArchetype, 
// which should be fine for most use cases.
type Archetype interface {
    // A unique identifier for an archetype
    GetHandle() ArchetypeHandle
    // Returns a sorted slice representing the Component layout for the archetype
    GetType() []ComponentHandle
    GetComponent(entity EntityHandle, componentHandle ComponentHandle) Component
    // Takes in a global EntityHandle and returns the local EntityHandle for the new entity
    CreateEntity(gobalHandle EntityHandle) EntityHandle
    // Takes in a local EntityHandle and removes that entity, also takes in a pointer to the 
    // ECS incase you need to modify anything there
    RemoveEntity(entity EntityHandle, ecs *ECS)
    // Local EntityHandle -> Gobal EntityHandle
    GetEntities() []EntityHandle
    // Get the edge of archetypes with/without the given component
    //GetEdge(component ComponentHandle) *ArchetypeEdge
    // Get the edges of this Archetype
    //GetEdges() []ArchetypeEdge
}

type GenericArchetype struct {
    Handle  ArchetypeHandle
    // Global ComponentHandle -> Local ComponentHandle
    CompMap map[ComponentHandle]ComponentHandle
    // should be sorted
    Type    []ComponentHandle
    // should have same order as Type
    // index of ComponentHandle in Type -> ComponentList of ComponentType
    Comps   []ComponentList
    // Local EntityHandle -> Gobal EntityHandle
    Ents    []EntityHandle
    // ComponentHandle -> ArchetypeEdge
    //Edges   []ArchetypeEdge
}

// Creates a GenericArchetype, should only be used by the ECS since the handle is defined by the ECS
func createArchetype(handle ArchetypeHandle, components []*ComponentType) Archetype {
    amap := make(map[ComponentHandle]ComponentHandle)
    atype := make([]ComponentHandle, len(components))
    comps := make([]ComponentList, len(components))
    for i := range components {
        ctype := components[i]
        amap[ctype.Handle] = ComponentHandle(i)
        atype[i] = ctype.Handle
        comps[i] = NewComponentList(*ctype)
    }

    return &GenericArchetype {
        Handle: handle,
        CompMap: amap,
        Type: atype,
        Comps: comps,
        Ents: make([]EntityHandle, 0, 1),
    }
}

func (a *GenericArchetype) GetHandle() ArchetypeHandle { return a.Handle }
func (a *GenericArchetype) GetType() []ComponentHandle { return a.Type }
func (a *GenericArchetype) GetEntities() []EntityHandle { return a.Ents }
func (a *GenericArchetype) GetComponent(entity EntityHandle, comp ComponentHandle) Component {
    c, has := a.CompMap[comp]
    if !has { return nil }
    return a.Comps[c].Get(entity)
}
func (a *GenericArchetype) CreateEntity(gobalHandle EntityHandle) EntityHandle {
    handle := EntityHandle(a.Comps[0].len)
    a.Ents = append(a.Ents, gobalHandle)
    for i := range len(a.Comps) {
        a.Comps[i].Add()
    }

    return handle
}
func (a *GenericArchetype) RemoveEntity(handle EntityHandle, ecs *ECS) {
    for i := range len(a.Comps) {
        a.Comps[i].Remove(handle)
    }

    a.Ents = slices.Delete(a.Ents, int(handle), int(handle) + 1)
    // update global entity records with the new local handle
    for i := uint64(handle); i < uint64(len(a.Ents)); i += 1 {
        global := a.Ents[i]
        record := &ecs.entities[global]
        record.Id -= 1
    }
}











