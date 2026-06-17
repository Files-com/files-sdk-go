package file

func (r *adaptiveManagerRegistry[K]) resetForTest() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.managers = nil
}
