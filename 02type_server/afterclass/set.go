package afterclass

type Set interface {
	Put(key string)
	Keys() []string
	Contains(key string) bool
	Remove(key string)
	// 如果之前已经有了，就返回旧的值，absent =false
	// 如果之前没有，就塞下去，返回 absent = true
	PutIfAbsent(key string) (old string, absent bool)
}

type Set1 map[string]bool

func (s Set1) Put(key string) {
	s[key] = true
}

func (s Set1) Keys() []string {
	//TODO implement me
	panic("implement me")
}

func (s Set1) Contains(key string) bool {
	//TODO implement me
	panic("implement me")
}

func (s Set1) Remove(key string) {
	//TODO implement me
	panic("implement me")
}

func (s Set1) PutIfAbsent(key string) (old string, absent bool) {
	//TODO implement me
	panic("implement me")
}

type MySet struct {
	M map[string]string
}

func (m *MySet) NewSet() {
	m.M = make(map[string]string)
}

func (m *MySet) Put(key string) {
	m.M[key] = ""
}

func (m MySet) Keys() []string {
	arr := make([]string, 0, len(m.M))
	for k, _ := range m.M {
		arr = append(arr, k)
	}
	return arr
}

func (m MySet) Contains(key string) bool {
	_, res := m.M[key]
	return res
}

func (m MySet) Remove(key string) {
	delete(m.M, key)
}

func (m MySet) PutIfAbsent(key string) (old string, absent bool) {
	_, res := m.M[key]
	if res {
		//contains
		return key, false
	} else {
		m.M[key] = ""
		return key, true
	}
}
