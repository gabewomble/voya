package server

import "fmt"

func (s *Server) background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				s.logger.LogError(nil, fmt.Errorf("panic: %v", err))
			}
		}()

		fn()
	}()
}
