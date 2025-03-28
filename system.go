package ecs




type System struct {
    // Tick interval between executions
    rate    uint
    // Number of ticks elapsed since last execution
    passed  uint
    Query   Query
    // Function to execute
    Fn      func(*ECS, []EntityHandle)
}

func (s *System) Check(ecs *ECS) {
    s.passed += 1
    if s.passed != s.rate { return }
    s.passed = 0
    s.Fn(ecs, s.Query.Apply(ecs))
}

func NewSystem(rate uint, query Query, fn func(*ECS, []EntityHandle)) System {
    return System {
        rate: rate,
        Query: query,
        Fn: fn,
    }
}
