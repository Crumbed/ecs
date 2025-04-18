package ecs

import "math"

//"fmt"

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

const InvalidHandle uint64 = math.MaxUint64
type EntityHandle uint64


var ComponentCount ComponentHandle = 0
type ECS struct {
    // ArchetypeHandle -> Archetype
    archetypes      []Archetype
    archetypeQuery  map[ComponentHash]ArchetypeHandle
    // EntityHandle -> ArchetypeRecord for the given entity
    entities        []Record
    // List of all removed entities
    RemovedEntities []EntityHandle

    // ComponentHandle -> Set[ArchetypeHandle] that all contain the given ComponentHandle
    componentIndex  []Set[ArchetypeHandle]
    // ComponentHandle -> ComponentType
    componentTypes  []ComponentType

    systems         []System
    //flags           map[string]EntityFlags
}

func NewECS() *ECS {
    return &ECS {
        archetypes: make([]Archetype, 0),
        archetypeQuery: make(map[ComponentHash]ArchetypeHandle),
        entities: make([]Record, 0),

        componentIndex: make([]Set[ArchetypeHandle], 0),
        componentTypes: make([]ComponentType, 0),

        //flags: make(map[string]EntityFlags),
    }
}

func (self *ECS) createArchetype(components []ComponentHandle, query ComponentHash) ArchetypeHandle {
    // create new archetype
    handle := ArchetypeHandle(len(self.archetypes))
    atype := make([]*ComponentType, len(components))
    for i, c := range components {
        atype[i] = &self.componentTypes[c]   
        self.componentIndex[c].Add(handle)
    }

    self.archetypes = append(self.archetypes, createArchetype(handle, atype))
    // update archetype query
    self.archetypeQuery[query] = handle

    return handle
}

func (self *ECS) HasComponent(e EntityHandle, component ComponentHandle) bool {
    record := &self.entities[e]
    archetype := record.Arch
    archset := self.componentIndex[component]
    return archset.Has(archetype.GetHandle())
}

func (self *ECS) GetComponent(e EntityHandle, component ComponentHandle) Component {
    record := self.entities[e]
    arch := record.Arch
    handle := arch.GetHandle()
    archset := self.componentIndex[component]
    if !archset.Has(handle) { return nil }
    return arch.GetComponent(record.Id, component)
}

func (self *ECS) GetComponentUnchecked(e EntityHandle, component ComponentHandle) Component {
    record := self.entities[e]
    return record.Arch.GetComponent(record.Id, component)
}

func (self *ECS) AddEntity(components ...ComponentHandle) EntityHandle {
    entity := EntityHandle(len(self.entities))
    c_query := CreateComponentHash(components...)
    // get archetype
    arch_h, exists := self.archetypeQuery[c_query]
    if !exists {
        arch_h = self.createArchetype(components, c_query)
    }
    arch := self.archetypes[arch_h]
    arch_entity := arch.CreateEntity(entity)
    
    self.entities = append(self.entities, Record {
        Arch: arch,
        Id: arch_entity,
    })

    return entity
}

func (self *ECS) RemoveEntity(handle EntityHandle) {
    record := &self.entities[handle]
    record.Arch.RemoveEntity(record.Id, self)
    record.Id = EntityHandle(InvalidHandle)
}

// Advance by 1 tick
func (self *ECS) Step() {
    for sysi := range len(self.systems) {
        sys := &self.systems[sysi]
        sys.Check(self)
    }
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
    if len(self.entities) != 0 {
        panic("Cannot register a new component after adding an entity!")
    }
    *handle = ComponentHandle(len(self.componentTypes))
    componentType := CreateComponentType(emptyComponent, *handle)
    /*
    components := NewComponentList(componentType)
    self.components = append(self.components, components)
    */
    self.componentTypes = append(self.componentTypes, componentType)
    self.componentIndex = append(self.componentIndex, NewSet[ArchetypeHandle]())
    ComponentCount += 1
    return self
}

func (self *ECS) WithArchetype(handle *ArchetypeHandle, arch Archetype, components ...ComponentHandle) *ECS {
    if len(self.entities) != 0 {
        panic("Cannot register a new archetype after adding an entity!")
    }
    *handle = ArchetypeHandle(len(self.archetypes))
    self.archetypes = append(self.archetypes, arch)
    for _, c := range components {
        self.componentIndex[c].Add(*handle)
    }
    query := CreateComponentHash(components...)
    self.archetypeQuery[query] = *handle

    return self
}

func (self *ECS) WithSystem(sys System) *ECS {
    self.systems = append(self.systems, sys)
    return self
}





