package utils

// Subtract returns a new slice that is a copy of input slice, subjected to the following
// procedure: for each element in the subtract slice, its first occurrence in the input
// slice is deleted.
func Subtract[V comparable](ivs, svs []V) []V {
	if ivs == nil || svs == nil {
		return ivs
	}
	ovs := Copy(ivs)
	for _, sv := range svs {
		ovs = Delete(sv, ovs)
	}
	return ovs
}

// Copy is simply a convenient combination of allocation and copying.
func Copy[V any](ivs []V) []V {
	if ivs == nil {
		return nil
	}
	ovs := make([]V, len(ivs))
	copy(ovs, ivs)
	return ovs
}

// Delete removes the first matching value of a slice.
func Delete[V comparable](dv V, ivs []V) []V {
	ovs := Copy(ivs)
	for i := range ovs {
		if ovs[i] == dv {
			ovs = append(ovs[:i], ovs[i+1:]...)
			return ovs
		}
	}
	return ovs
}
