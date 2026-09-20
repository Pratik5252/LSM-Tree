package memtable


type Entry struct {
	Value []byte
	Tombstone bool
}
type MemTable struct {
	data map[string]Entry
}

// Create new MemTable instance
func New() *MemTable {
	return &MemTable{
		data: make(map[string]Entry),
	}
}

// Put adds key-vlaue pair to MemTable
func (m *MemTable) Put(key string, value []byte) {
	m.data[key] = Entry{value,false}
}

// Gets value by key and return value in []byte adn a bool as status
func (m *MemTable) Get(key string) (Entry, bool) {
	entry, ok := m.data[key]
	
	if !ok {
		return Entry{}, false
	}

	return entry, true
}

// Delete a key-value pair
func (m *MemTable) Delete(key string) {
	m.data[key] = Entry{nil,true}
}

func (m *MemTable) Count() int {
	return len(m.data)
}

func (m *MemTable) Entries() map[string]Entry {
	return m.data
}

func (m *MemTable) Clear() {
	clear(m.data)
}