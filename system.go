package ecs





type System struct {
    rate    uint
    passed  uint
    query   ComponentHash
    Fn      func(*ECS)
}

func (s *System) Check(ecs *ECS) {
    s.passed += 1
    if s.passed != s.rate { return }
    s.passed = 0
    s.Fn(ecs)
}

func NewSystem(rate uint, query ComponentHash, fn func(*ECS)) System {
    return System {
        rate: rate,
        query: query,
        Fn: fn,
    }
}
