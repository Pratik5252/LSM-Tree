package memtable


type MemTable struct {
	data map[string][]byte
}

// Create new MemTable instance
func New() *MemTable {
	return &MemTable{
		data: make(map[string][]byte),
	}
}

// Put adds key-vlaue pair to MemTable
func (m *MemTable) Put(key string, value []byte) {
	m.data[key] = value
}

// Gets value by key and return value in []byte adn a bool as status
func (m *MemTable) Get(key string) ([]byte, bool) {
	value, ok := m.data[key]

	if !ok {
		return nil, ok
	}

	return value, ok
}

// Delete a key-value pair
func (m *MemTable) Delete(key string) {
	delete(m.data, key)
}

func (m *MemTable) Count() int {
	return len(m.data)
}

func (m *MemTable) Entries() map[string][]byte {
	return m.data
}