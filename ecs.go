package ecs

import (
)


type ECS struct {
	components      []ComponentList
    componentTypes  []ComponentType
}

func NewECS() *ECS {
    return &ECS {
        components: make([]ComponentList, 0),
        componentTypes: make([]ComponentType, 0),
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
    *handle = ComponentHandle(len(self.components))
    componentType := CreateComponentType(emptyComponent)
    components := NewComponentList(componentType)
    self.components = append(self.components, components)
    self.componentTypes = append(self.componentTypes, componentType)
    return self
}





