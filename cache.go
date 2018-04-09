package main

type CacheKey struct {
	name string
	key  string
}

type Cache struct {
	Store
	data map[CacheKey]Item
}

func NewCache(s Store) *Cache {
	c := &Cache{}
	c.Store = s
	c.data = make(map[CacheKey]Item)
	return c
}

func (c *Cache) Get(name, key string) (string, error) {
	item, err := c.GetItem(name, key)
	if err != nil {
		return "", err
	}
	return item.Find(key)
}

func (c *Cache) GetItem(name, key string) (*Item, error) {
	ck := CacheKey{name, key}
	if item, ok := c.data[ck]; ok {
		return &item, nil
	}

	item, err := c.Store.GetItem(name, key)
	if err != nil {
		return &Item{}, err
	}
	c.data[ck] = *item
	return item, nil
}
