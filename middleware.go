package gogator

import "context"

func MiddlewareLoggedIn(hl handlerLoggedIn) handler {
	return func(s *state, c command) error {
		u, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}
		return hl(s, c, u)
	}
}
