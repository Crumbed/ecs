package ecs

import (
    "math"
)

type Set[T comparable] map[T]struct{}
func NewSet[T comparable](elements ...T) Set[T] {
    s := make(Set[T])
    for _, e := range elements {
        s.Add(e)
    }
    return s
}

func (s Set[T]) Add(e T) { s[e] = struct{}{} }
func (s Set[T]) Has(e T) bool {
    _, has := s[e]
    return has
}

type ArchetypeHandle uint64
type Archetype struct {
    Handle  ArchetypeHandle
    Type    []ComponentHandle
    //Set     Set[ComponentHandle]
}

type EntityHandle uint64
const InvalidEntityHandle EntityHandle = math.MaxUint64
type Entity struct {
    Id      EntityHandle
    Mask    uint64
}


type ECS struct {
    // ArchetypeHandle -> Archetype
    archetypes      []Archetype
    archetypeQuery  map[ComponentQuery]ArchetypeHandle
    // ComponentHandle -> Set[ArchetypeHandle] that all contain the given ComponentHandle
    componentIndex  []Set[ArchetypeHandle]
    // EntityHandle -> *Archetype of given entity
    entities        []*Archetype
    // ComponentHandle -> ComponentList of proper ComponentType
	components      []ComponentList
    // ComponentHandle -> ComponentType
    componentTypes  []ComponentType
}

func NewECS() *ECS {
    return &ECS {
        archetypes: make([]Archetype, 0),
        archetypeQuery: make(map[ComponentQuery]ArchetypeHandle),
        componentIndex: make([]Set[ArchetypeHandle], 0),
        entities: make([]*Archetype, 0),
        components: make([]ComponentList, 0),
        componentTypes: make([]ComponentType, 0),
    }
}

func (self *ECS) HasComponent(e EntityHandle, component ComponentHandle) bool {
    archetype := self.entities[e]
    archset := self.componentIndex[component]
    return archset.Has(archetype.Handle)
}

// `handle *ComponentHandle` is set to a value, NEVER CHANGE THIS VALUE AFTER REGISTERING A COMPONENT.
// `components ComponentList` represents the list of entity components.
//  EXAMPLE CALL:
// 
//  var HealthHandle ecs.ComponentHandle
//  type Health struct {...}
//  func (h *Health) GetComponentHandle() ecs.ComponentHandle { 
//      return HealthHandle 
//  }
// 
//  func main() {
//      ecs := ecs.NewECS().
//          WithComponentType(&HealthHandle, &Health{})
//  }
func (self *ECS) WithComponentType(handle *ComponentHandle, emptyComponent Component) *ECS {
    *handle = ComponentHandle(len(self.components))
    componentType := CreateComponentType(emptyComponent)
    components := NewComponentList(componentType)
    self.components = append(self.components, components)
    self.componentTypes = append(self.componentTypes, componentType)
    return self
}





