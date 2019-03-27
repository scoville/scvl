package engine

// AuthCodeURL returns the AuthCodeURL of google
func (e *Engine) AuthCodeURL(state string) string {
	return e.googleClient.AuthCodeURL(state)
}
