package ecs

import (
	//"fmt"
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
    Type    []ComponentHandle // maybe this should be a ComponentQuery
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
    // EntityHandle -> ArchetypeHandle of given entity
    entities        []ArchetypeHandle
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
        entities: make([]ArchetypeHandle, 0),
        components: make([]ComponentList, 0),
        componentTypes: make([]ComponentType, 0),
    }
}

func (self *ECS) createArchetype(components []ComponentHandle, query ComponentQuery) ArchetypeHandle {
    // create new archetype
    handle := ArchetypeHandle(len(self.archetypes))
    self.archetypes = append(self.archetypes, Archetype {
        Type: components,
        Handle: handle,
    })

    // update component index with new archetype
    for _, c := range components {
        self.componentIndex[c].Add(handle)
    }

    // update archetype query
    self.archetypeQuery[query] = handle

    return handle
}

func (self *ECS) HasComponent(e EntityHandle, component ComponentHandle) bool {
    ahandle := self.entities[e]
    archetype := self.archetypes[ahandle]
    archset := self.componentIndex[component]
    return archset.Has(archetype.Handle)
}

func (self *ECS) GetComponent(e EntityHandle, component ComponentHandle) Component {
    comp_list := self.components[component]
    return comp_list.Get(uint64(e))
}

func (self *ECS) AddEntity(components ...ComponentHandle) EntityHandle {
    entity := EntityHandle(len(self.entities))
    c_query := CreateQuery(components...)
    // get archetype
    arch_h, exists := self.archetypeQuery[c_query]
    if !exists {
        arch_h = self.createArchetype(components, c_query)
    }
    
    self.entities = append(self.entities, arch_h)
    for i := range self.components {
        self.components[i].Add()
    }
    return entity
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
    self.componentIndex = append(self.componentIndex, NewSet[ArchetypeHandle]())
    return self
}





