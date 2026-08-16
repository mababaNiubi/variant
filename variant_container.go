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
		mp, _ := c.mapVariant()
		return len(mp)
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
		mp, ok := c.mapVariant()
		if !ok {
			return
		}
		for key, value := range mp {
			if !f(key, value) {
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
		mp, ok := c.mapVariant()
		if !ok {
			return
		}
		index := 0
		for _, value := range mp {
			if !f(index, value) {
				return
			}
			index++
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
		mp, ok := c.mapVariant()
		if !ok {
			return nil
		}
		delete(mp, key)
		c.complexValue = mp
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
		switch m := c.complexValue.(type) {
		case map[string]Variant:
			v, has := m[key]
			return v, has
		case map[string]any:
			// Raw structure: wrap only the requested value, never convert the
			// whole map (this runs per point during query condition evaluation).
			val, has := m[key]
			if !has {
				return NewEmpty(), false
			}
			return NewRawValue(val), true
		}
	}
	return NewEmpty(), false
}

func (c *Variant) MapSet(key string, value any) {
	switch c.variantType {
	case TypeMap:
		mp, ok := c.mapVariant()
		if !ok {
			return
		}
		mp[key] = New(value)
		c.complexValue = mp
	case TypeEmpty:
		mp := make(map[string]Variant)
		mp[key] = New(value)
		c.complexValue = mp
		c.variantType = TypeMap
	default:
		return
	}
}
