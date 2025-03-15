package ecs

import (
    "testing"
)

var HealthHandle ComponentHandle
type Health struct {
    hp uint64
}
func (h *Health) GetComponentHandle() ComponentHandle { return HealthHandle }
func TestWithComponentType(t *testing.T) {
    ecs := NewECS()
    ecs.WithComponentType(&HealthHandle, &Health{})
}
