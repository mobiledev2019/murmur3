package murmur3

import (
	"log"
)

type ByteBuffer struct {
	HB              []byte
	Offset          int
	Mark            int
	Position        int
	Limit           int
	Capacity        int
	BigEndian       bool
	NativeByteOrder bool
}

func allocateByteBuffer(bufferSize int) *ByteBuffer {
	return &ByteBuffer{
		HB:              make([]byte, bufferSize),
		Offset:          0,
		Mark:            -1,
		Position:        0,
		Limit:           bufferSize,
		Capacity:        bufferSize,
		NativeByteOrder: true,
	}
}

func InitMurmurByteBuffer(chunkSize int) *MurmurByteBuffer {
	mBuff := &MurmurByteBuffer{}
	mBuff.chunkSize = chunkSize
	mBuff.bufferSize = chunkSize
	mBuff.buffer = allocateByteBuffer(mBuff.bufferSize + 7)
	return mBuff
}

func (buf *ByteBuffer) putCharL(index int, val byte) {
	buf.HB[index] = val
	buf.HB[index+1] = val >> 8
}

func (buf *ByteBuffer) remaining() int {
	return buf.Limit - buf.Position
}

func (buf *ByteBuffer) flip() {
	buf.Limit = buf.Position
	buf.Position = 0
	buf.Mark = -1
}

func (buf *ByteBuffer) getLongB(index int) int64 {
	return MakeLong(Get(buf.HB, index), Get(buf.HB, index+1), Get(buf.HB, index+2), Get(buf.HB, index+3), Get(buf.HB, index+4), Get(buf.HB, index+5), Get(buf.HB, index+6), Get(buf.HB, index+7))
}

func (buf *ByteBuffer) getLongL(index int) int64 {
	return MakeLong(Get(buf.HB, index+7), Get(buf.HB, index+6), Get(buf.HB, index+5), Get(buf.HB, index+4), Get(buf.HB, index+3), Get(buf.HB, index+2), Get(buf.HB, index+1), Get(buf.HB, index))
}

func (buf *ByteBuffer) nextGetIndex(var1 int) int {
	if buf.Limit-buf.Position < var1 {
		log.Println("[E] Buffer Underflow Exception")
	} else {
		var2 := buf.Position
		buf.Position += var1
		return var2
	}
	return -1
}

func (buf *ByteBuffer) getLong() int64 {
	if buf.BigEndian {
		return buf.getLongB(buf.ix(buf.nextGetIndex(8)))
	} else {
		return buf.getLongL(buf.ix(buf.nextGetIndex(8)))
	}
}

func (buf *ByteBuffer) capacity() int {
	return buf.Capacity
}

func (buf *ByteBuffer) discardMark() {
	buf.Mark = -1
}

func (buf *ByteBuffer) compact() {
	buf.position(buf.remaining())
	buf.limit(buf.capacity())
	buf.discardMark()
}

func (buf *ByteBuffer) get(index int) byte {
	return buf.HB[index]
}

func (buf *ByteBuffer) position(var1 int) {
	if var1 <= buf.Limit && var1 >= 0 {
		buf.Position = var1
		if buf.Mark > buf.Position {
			buf.Mark = -1
		}
	} else {
		log.Println("[E] Illegal Argument Exception")
	}
}

func (buf *ByteBuffer) limit(chunkSize int) {
	if chunkSize <= buf.Capacity && chunkSize >= 0 {
		buf.Limit = chunkSize
		if buf.Position > buf.Limit {
			buf.Position = buf.Limit
		}

		if buf.Mark > buf.Limit {
			buf.Mark = -1
		}
	} else {
		log.Println("[E] Illegal Argument Exception")
	}
}

func (buf *ByteBuffer) ix(nextPutIndex int) int {
	return buf.Offset + nextPutIndex
}

func (buf *ByteBuffer) nextPutIndex(var1 int) int {
	if buf.Limit-buf.Position < var1 {
		log.Println("[E] Buffer Overflow Exception!")
		return -1
	} else {
		var2 := buf.Position
		buf.Position += var1
		return var2
	}
}

func (buf *ByteBuffer) putLongB(index int, val int64) {
	buf.HB[index] = Long7(val)
	buf.HB[index+1] = Long6(val)
	buf.HB[index+2] = Long5(val)
	buf.HB[index+3] = Long4(val)
	buf.HB[index+4] = Long3(val)
	buf.HB[index+5] = Long2(val)
	buf.HB[index+6] = Long1(val)
	buf.HB[index+7] = Long0(val)
}

func (buf *ByteBuffer) putLongL(index int, val int64) {
	buf.HB[index+7] = Long7(val)
	buf.HB[index+6] = Long6(val)
	buf.HB[index+5] = Long5(val)
	buf.HB[index+4] = Long4(val)
	buf.HB[index+3] = Long3(val)
	buf.HB[index+2] = Long2(val)
	buf.HB[index+1] = Long1(val)
	buf.HB[index] = Long0(val)
}

func (buf *ByteBuffer) putLong(val int64) {
	if buf.BigEndian {
		buf.putLongB(buf.ix(buf.nextPutIndex(8)), val)
	} else {
		buf.putLongL(buf.ix(buf.nextPutIndex(8)), val)
	}
}

func (buf *ByteBuffer) checkIndex(index int) int {
	if index >= 0 && index < buf.Limit {
		return index
	} else {
		log.Println(" [E] Array Index Out Of Bounds Exception")
	}
	return -1
}

/*
	The first 4 bytes to integer 32 bit
*/
func (buf *ByteBuffer) AsInt() int {
	return  int(int32(buf.HB[0]) & 255 | (int32(buf.HB[1]) & 255) << 8 | (int32(buf.HB[2]) & 255) << 16 | (int32(buf.HB[3]) & 255) << 24)
}

/*
	Trả về mảng 4 bytes đầu tiên, theo thứ tự ngược lại, dùng để ghép key sử dụng get row trong Hbase
*/
func (buf *ByteBuffer) AsBytes() []byte {
	return []byte{buf.HB[3], buf.HB[2], buf.HB[1], buf.HB[0]}
}

/*
	Trả về 16 bytes, chính là H1, H2 ở dạng bytes
*/
func (buf *ByteBuffer) ToBytes() []byte {
	return buf.HB
}

/*
	Trả về chuỗi hash dài 32 bytes được tạo từ 16 bytes ban đầu của H1 và H2
*/
func (buf *ByteBuffer) ToString() string {
	return ToString(buf.HB)
}