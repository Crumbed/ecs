package ecs

import (
	"fmt"
	"testing"
)

var HealthHandle ComponentHandle
type Health struct {
    hp uint64
}
func (h *Health) GetComponentHandle() ComponentHandle { return HealthHandle }

var PosHandle ComponentHandle
type Pos struct {
    x, y uint64
}
func (p *Pos) GetComponentHandle() ComponentHandle { return PosHandle }

func TestWithComponentType(t *testing.T) {
    fmt.Println("\n[Testing component types]")
    ecs := NewECS()
    ecs.WithComponentType(&HealthHandle, &Health{})
    ecs.WithComponentType(&PosHandle, &Pos{})
    fmt.Println("HealthHandle:", HealthHandle)
    fmt.Println("PosHandle:", PosHandle)
}

func TestArchetypes(t *testing.T) {
    fmt.Println("\n[Testing archetypes]")
    ecs := NewECS().
        WithComponentType(&HealthHandle, &Health{}).
        WithComponentType(&PosHandle, &Pos{})

    ecs.AddEntity(HealthHandle, PosHandle)
    ecs.AddEntity(HealthHandle, PosHandle)
    ecs.AddEntity(HealthHandle)
    ecs.AddEntity(PosHandle)

    for _, a := range ecs.archetypes {
        fmt.Println(a)
    }
}


func TestComponents(t *testing.T) {
    fmt.Println("\n[Testing components]")
    ecs := NewECS().
        WithComponentType(&HealthHandle, &Health{}).
        WithComponentType(&PosHandle, &Pos{})

    // create
    e1 := ecs.AddEntity(HealthHandle, PosHandle)
    e2 := ecs.AddEntity(HealthHandle, PosHandle)

    // initialize
    health := ecs.GetComponent(e1, HealthHandle).(*Health)
    pos := ecs.GetComponent(e1, PosHandle).(*Pos)
    health.hp = 10
    pos.x = 1
    pos.y = 2

    health = ecs.GetComponent(e2, HealthHandle).(*Health)
    pos = ecs.GetComponent(e2, PosHandle).(*Pos)
    health.hp = 100
    pos.x = 10
    pos.y = 10

    // check
    e1h := ecs.GetComponent(e1, HealthHandle).(*Health)
    e1p := ecs.GetComponent(e1, PosHandle).(*Pos)
    if e1h.hp != 10 || (e1p.x != 1 || e1p.y != 2) {
        t.Errorf("Expected { health: 10, x: 1, y: 2 } but found { health: %d, x: %d, y: %d }", e1h.hp, e1p.x, e1p.y)
    }
    e2h := ecs.GetComponent(e2, HealthHandle).(*Health)
    e2p := ecs.GetComponent(e2, PosHandle).(*Pos)
    if e2h.hp != 100 || (e2p.x != 10 || e2p.y != 10) {
        t.Errorf("Expected { health: 100, x: 10, y: 10 } but found { health: %d, x: %d, y: %d }", e2h.hp, e2p.x, e2p.y)
    }
}

func TestQueriesAndSystems(t *testing.T) {
    fmt.Println("\n[Testing queries & systems]")
    ecs := NewECS().
        WithComponentType(&HealthHandle, &Health{}).
        WithComponentType(&PosHandle, &Pos{})
    
    ecs.AddEntity(HealthHandle, PosHandle)
    ecs.AddEntity(HealthHandle, PosHandle)
    ecs.AddEntity(HealthHandle)
    ecs.AddEntity(PosHandle)
}














