package mock

// GetInterface is a helper to pull interfaces off interface channels
func GetInterface(channel chan interface{}) interface{} {
	select {
	case chanObj := <-channel:
		return chanObj
	default:
		return nil
	}
}
