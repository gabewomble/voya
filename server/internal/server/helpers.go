package server

func (s *Server) background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				s.log.LogError(nil, "Failed to recover background call", nil, "error", err)
			}
		}()

		fn()
	}()
}
