package backend

import (
	"context"
	"strconv"
	"strings"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

var _ Expressions = (*impl)(nil)

type impl struct {
	logger *logger.Logger
	repo   interfaces.Repository
}

func NewBackend(logger *logger.Logger, repo interfaces.Repository) *impl {
	return &impl{
		logger: logger,
		repo:   repo,
	}
}

func (s *impl) Process(ctx context.Context, id, exp string) (*int, error) {
	// mark the expression in progress in the backend data store
	err := s.repo.MarkInProgress(ctx, id)

	if err != nil {
		return nil, err
	}

	val, err := s.process(exp)

	if err != nil {
		msg := err.Error()
		err := s.repo.MarkFailed(ctx, id, &msg)
		return nil, err
	}

	s.repo.MarkCompleted(ctx, id, &val)

	return &val, nil
}

func (s *impl) process(exp string) (int, error) {
	tokens := strings.Split(exp, ",")
	return s.evalRPN(tokens)
}

func (s *impl) evalRPN(tokens []string) (int, error) {
	operations := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}
	stack := []int{}
	for _, token := range tokens {
		if operation, exists := operations[token]; exists {
			num1 := stack[len(stack)-2]
			num2 := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			result := operation(num1, num2)
			stack = append(stack, result)
		} else {
			num, _ := strconv.Atoi(token)
			stack = append(stack, num)
		}
	}
	return stack[0], nil
}
