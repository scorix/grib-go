package bitio

import (
	"fmt"
	"io"
)

// readBits reads n bits from byte slice starting at given bit offset
func ReadBits(data []byte, offset, n uint8) (uint64, error) {
	if n > 64 {
		return 0, fmt.Errorf("cannot read more than 64 bits at once")
	}

	var result uint64
	currentByte := offset / 8
	bitPos := offset % 8

	bitsRemaining := n
	for bitsRemaining > 0 {
		if int(currentByte) >= len(data) {
			return 0, io.EOF
		}

		// How many bits we can read from current byte
		bitsAvailable := 8 - bitPos
		bitsToRead := bitsRemaining
		if bitsToRead > bitsAvailable {
			bitsToRead = bitsAvailable
		}

		// Create mask and read bits
		mask := byte((1 << bitsToRead) - 1)
		bits := (data[currentByte] >> (8 - bitPos - bitsToRead)) & mask

		// Add to result
		result = (result << bitsToRead) | uint64(bits)

		bitsRemaining -= bitsToRead
		currentByte++
		bitPos = 0
	}

	return result, nil
}
