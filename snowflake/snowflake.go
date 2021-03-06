// Package snowflake is a distributed unique ID generator based on
// Twitter's snowflake service
package snowflake

import (
	"time"
)

var workerIDBits uint64 = 5
var datacenterIDBits uint64 = 5
var sequenceBits uint64 = 12

var workerIDShift uint64 = sequenceBits
var datacenterIDShift uint64 = sequenceBits + workerIDBits
var timestampLeftShift uint64 = sequenceBits + workerIDBits + datacenterIDBits
var sequenceMask int64 = int64(1 << sequenceBits) - 1

var twepoch int64 = 1288834974657

// Snowflake struct is responsible for generating unique IDs
type Snowflake struct {
	workerID uint64
	datacenterID uint64
	sequenceNumber uint64
	lastTimestamp int64
}

// NewSnowflake is the constructor for snowflakes. It is considered a
// programming error to create a Snowflake by any other means.
func NewSnowflake(datacenterID uint64, workerID uint64) *Snowflake {
	s := new(Snowflake)
	s.datacenterID = datacenterID
	s.workerID = workerID
	s.sequenceNumber = 0
	s.lastTimestamp = 0
	return s
}

// NextID returns a unique id value generated by the
// snowflake. Calling NextId multiple times on the same snowflake will
// return unique, monotonically increasing IDs.
func (s *Snowflake) NextID() uint64 {

	timestamp := timeGen()

	if (timestamp < s.lastTimestamp) {
		panic("clock is going backward!")
	}

	if (s.lastTimestamp == timestamp) {
		s.sequenceNumber = uint64((s.sequenceNumber + 1) & uint64(sequenceMask))
		if (s.sequenceNumber == 0) {
			timestamp = tilNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequenceNumber = 0
	}

	s.lastTimestamp = timestamp

	return uint64(((timestamp - twepoch) << timestampLeftShift)) | 
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequenceNumber
}


// DatacenterID returns the snowflake's configured datacenterID
func (s *Snowflake) DatacenterID() uint64 {
	return s.datacenterID
}

// WorkerID returns the snowflake's configured workerID
func (s *Snowflake) WorkerID() uint64 {
	return s.workerID
}


//// package private

func timeGen() int64 {
	nanos := time.Now().UnixNano()
	millis := nanos / 1000000
	return millis
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}
