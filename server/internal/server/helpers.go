package server

import "strconv"

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

func parseStringToInt32(s string, defaultValue int32) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int32(i)
}
