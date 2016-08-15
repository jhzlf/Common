// IDGenerator.go
package Common

//type ID uint64
//type ID32 uint32

type IDGenerator struct {
	incID   uint64
	incID32 uint32
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		incID:   0,
		incID32: 0,
	}
}

//start id is 1.
func (g *IDGenerator) NewID() uint64 {
	g.incID++
	return g.incID
}

func (g *IDGenerator) NewID32() uint32 {
	g.incID32++
	return g.incID32
}

func (g *IDGenerator) Reset() {
	g.incID = 0
	g.incID32 = 0
}
