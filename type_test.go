package goapidoc

import (
	"log"
	"testing"
)

func TestParseInnerType(t *testing.T) {
	str := "aType<T1<>, T2<TT1<integer[]>>, T3<TT1, TT2>, T4>[][]"
	res := parseInnerType(str)
	obj := res.OutItems.Type.OutItems.Type
	log.Println(res.Name, res.OutItems)
	log.Println(res.OutItems.Type.Name, res.OutItems.Type.OutItems)
	log.Println(obj.Name, obj.OutSchema.Type, obj.OutSchema.Generic)

	g0 := obj.OutSchema.Generic[0]   // T1<>
	g1 := obj.OutSchema.Generic[1]   // T2<TT1<integer[]>>
	g10 := g1.OutSchema.Generic[0]   // TT1<integer[]>
	g100 := g10.OutSchema.Generic[0] // integer[]
	g1000 := g100.OutItems.Type      // integer
	g2 := obj.OutSchema.Generic[2]   // T3<TT1, TT2>
	g20 := g2.OutSchema.Generic[0]   // TT1
	g21 := g2.OutSchema.Generic[1]   // TT2
	g3 := obj.OutSchema.Generic[3]   // T4
	log.Println(g0.Name, g0.OutSchema.Type)
	log.Println(g1.Name, g1.OutSchema.Type, g10.OutSchema.Type, g100.Name, g1000.Name)
	log.Println(g2.Name, g2.OutSchema.Type, g20.OutSchema.Type, g21.OutSchema.Type)
	log.Println(g3.Name, g3.OutSchema.Type)
}
