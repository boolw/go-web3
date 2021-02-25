package abi

func (t *Type) Clone() *Type {
	if t == nil {
		return nil
	}
	item := new(Type)
	item.kind = t.kind
	item.size = t.size
	item.elem = t.elem.Clone()
	item.raw = t.raw
	if t.tuple != nil {
		item.tuple = make([]*TupleElem, len(t.tuple))
		for k, v := range t.tuple {
			item.tuple[k] = &TupleElem{
				//Name: v.Name,
				Elem:    v.Elem.Clone(),
				Indexed: v.Indexed,
			}
		}
	}
	item.t = t.t
	return item
}

func (e *Event) Clone() *Event {
	if e == nil {
		return nil
	}
	item := new(Event)
	item.Name = e.Name
	item.Anonymous = e.Anonymous
	item.Inputs = e.Inputs.Clone()
	item.id = e.id
	return item
}

func (m *Method) Clone() *Method {
	if m == nil {
		return nil
	}
	item := new(Method)
	item.Name = m.Name
	item.Const = m.Const
	item.Inputs = m.Inputs.Clone()
	item.Outputs = m.Outputs.Clone()
	item.id = m.id
	return item
}
