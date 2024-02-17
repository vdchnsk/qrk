package evaluator

type HashMap struct {
	store []string
}

func strToInt(str string, tableSize int) int {
	orbitraryInitialValue := 20
	orbitraryHashMultValue := 13

	hash := orbitraryInitialValue

	for _, letter := range str {
		hash *= orbitraryHashMultValue * int(letter) % tableSize
	}

	return hash
}

func NewHashMap() HashMap {
	hm := HashMap{
		store: make([]string, 100),
	}
	return hm
}

func (hm *HashMap) GetItem(key string) string {
	index := strToInt(key, len(hm.store))
	return hm.store[index]
}

func (hm *HashMap) SetItem(key, value string) {
	index := strToInt(key, len(hm.store))
	hm.store[index] = value
}
