package generics

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"unsafe"
)

// Hasher is an interface primarily for structs to provide a hash of their contents
type Hasher interface {
	Hash() uint32
}

// HashBytes returns a 32 bit unsigned integer hash of the passed byte slice
func HashBytes(b []byte) uint32 {
	hash := fnv.New32a()
	_, err := hash.Write(b)
	if err != nil {
		return 0
	}

	return hash.Sum32()
}

// Hash returns a 32 bit unsigned integer hash for any value passed in
func Hash[T comparable](value T) uint32 {
	hash := fnv.New32a()
	buf := new(bytes.Buffer)

	switch v := any(value).(type) {
	case Hasher:
		return v.Hash()

	case int:
		err := binary.Write(buf, binary.LittleEndian, int64(v))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	case *int:
		err := binary.Write(buf, binary.LittleEndian, uint64(uintptr(unsafe.Pointer(v))))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	case uint:
		err := binary.Write(buf, binary.LittleEndian, uint64(v))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	case *uint:
		err := binary.Write(buf, binary.LittleEndian, uint64(uintptr(unsafe.Pointer(v))))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	case uintptr:
		err := binary.Write(buf, binary.LittleEndian, uint64(v))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	case string:
		_, err := hash.Write([]byte(v))
		if err != nil {
			return 0
		}

	case *string:
		err := binary.Write(buf, binary.LittleEndian, uint64(uintptr(unsafe.Pointer(v))))
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}

	default:
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0
		}

		_, err = hash.Write(buf.Bytes())
		if err != nil {
			return 0
		}
	}

	return hash.Sum32()
}
