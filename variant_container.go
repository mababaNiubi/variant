package variant

import (
	"errors"
	"strconv"
)

func (c Variant) IsContainer() bool {
	if c.variantType == TypeMap || c.variantType == TypeList {
		return true
	}
	return false
}

func (c *Variant) Len() int {
	switch c.variantType {
	case TypeMap:
		sp, _ := c.StructPairs()
		return sp.len()
	case TypeList:
		return len(c.complexValue.([]Variant))
	case TypeString:
		return len(c.AsString())
	default:
		return 0
	}
}

func (c *Variant) Add(value any) (*Variant, error) {
	switch c.variantType {
	case TypeList:
		list, ok := c.complexValue.([]Variant)
		if !ok {
			return c, nil
		}
		list = append(list, New(value))
		c.complexValue = list
	case TypeEmpty:
		list := make([]Variant, 1)
		list[0] = New(value)
		c.complexValue = list
		c.variantType = TypeList
	case TypeString:
		str, ok := value.(string)
		if ok {
			c.AddString(str)
		}
	default:
		return nil, errors.New(errUnsupportedType)
	}
	return c, nil
}

func (c *Variant) Range(f func(key string, value Variant) bool) {
	switch c.variantType {
	case TypeMap:
		sp, _ := c.StructPairs()
		for i := range sp.keys {
			if !f(sp.keys[i], sp.vals[i]) {
				return
			}
		}
	case TypeList:
		list, ok := c.complexValue.([]Variant)
		if !ok {
			return
		}
		for i, variant := range list {
			if !f(strconv.Itoa(i), variant) {
				return
			}
		}
	default:
		return
	}
}

func (c *Variant) RangeByIndex(f func(index int, value Variant) bool) {
	switch c.variantType {
	case TypeMap:
		sp, _ := c.StructPairs()
		for i := range sp.vals {
			if !f(i, sp.vals[i]) {
				return
			}
		}
	case TypeList:
		list, ok := c.complexValue.([]Variant)
		if !ok {
			return
		}
		for i, variant := range list {
			if !f(i, variant) {
				return
			}
		}
	default:
		return
	}
}

func (c *Variant) Remove(i any) error {
	switch c.variantType {
	case TypeList:
		index, err := New(i).AsInt()
		if err != nil {
			return err
		}
		list, ok := c.complexValue.([]Variant)
		if !ok || len(list) <= index {
			return nil
		}
		c.complexValue = append(list[:index], list[index+1:]...)
	case TypeMap:
		key := New(i).AsString()
		sp, ok := c.StructPairs()
		if !ok {
			return nil
		}
		for j, k := range sp.keys {
			if k == key {
				sp.keys = append(sp.keys[:j], sp.keys[j+1:]...)
				sp.vals = append(sp.vals[:j], sp.vals[j+1:]...)
				break
			}
		}
		c.complexValue = sp
	default:
		return errors.New(errUnsupportedType)
	}
	return nil
}

func (c *Variant) Get(i any) (Variant, bool) {
	switch c.variantType {
	case TypeList:
		index, err := New(i).AsInt()
		if err != nil {
			return NewEmpty(), false
		}
		return c.ListGet(index)
	case TypeMap:
		key := New(i).AsString()
		return c.MapGet(key)
	default:
		return NewEmpty(), false
	}
}

func (c *Variant) Set(i any, value any) {
	switch c.variantType {
	case TypeList:
		index, err := New(i).AsInt()
		if err == nil {
			c.ListSet(index, value)
		}
	case TypeMap:
		key := New(i).AsString()
		c.MapSet(key, value)
	default:
	}
}

func (c *Variant) ListGet(index int) (Variant, bool) {
	if index < 0 {
		return NewEmpty(), false
	}
	if c.variantType == TypeList {
		list, ok := c.complexValue.([]Variant)
		if !ok {
			return NewEmpty(), false
		}
		if len(list) > index {
			return list[index], true
		}
	} else if index == 0 {
		return *c, false
	}
	return NewEmpty(), false
}

func (c *Variant) ListSet(index int, value any) {
	if index < 0 {
		return
	}
	if c.variantType == TypeList {
		list, ok := c.complexValue.([]Variant)
		if !ok {
			return
		}
		if len(list) <= index {
			return
		}
		list[index] = New(value)
		c.complexValue = list
	}
	return
}

func (c *Variant) MapGet(key string) (Variant, bool) {
	if c.variantType == TypeMap {
		sp, ok := c.StructPairs()
		if !ok {
			return NewEmpty(), false
		}
		return sp.get(key)
	}
	return NewEmpty(), false
}

func (c *Variant) MapSet(key string, value any) {
	switch c.variantType {
	case TypeMap:
		sp, ok := c.StructPairs()
		if !ok {
			return
		}
		nv := New(value)
		for j, k := range sp.keys {
			if k == key {
				sp.vals[j] = nv
				return
			}
		}
		sp.keys = append(sp.keys, key)
		sp.vals = append(sp.vals, nv)
		c.complexValue = sp
	case TypeEmpty:
		sp := &structPairs{keys: []string{key}, vals: []Variant{New(value)}}
		c.complexValue = sp
		c.variantType = TypeMap
	default:
		return
	}
}
