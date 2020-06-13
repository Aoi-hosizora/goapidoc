package goapidoc

import (
	"log"
	"testing"
)

func TestParseInnerType(t *testing.T) {
	str := "aType<T1<>, T2<TT1<integer[]>>, T3<TT1, TT2>, T4>[][]"
	res := parseInnerType(str)
	obj := res.OutArray.Type.OutArray.Type
	log.Println(res.Name, res.OutArray)
	log.Println(res.OutArray.Type.Name, res.OutArray.Type.OutArray)
	log.Println(obj.Name, obj.OutObject.Type, obj.OutObject.Generic)

	g0 := obj.OutObject.Generic[0]   // T1<>
	g1 := obj.OutObject.Generic[1]   // T2<TT1<integer[]>>
	g10 := g1.OutObject.Generic[0]   // TT1<integer[]>
	g100 := g10.OutObject.Generic[0] // integer[]
	g1000 := g100.OutArray.Type      // integer
	g2 := obj.OutObject.Generic[2]   // T3<TT1, TT2>
	g20 := g2.OutObject.Generic[0]   // TT1
	g21 := g2.OutObject.Generic[1]   // TT2
	g3 := obj.OutObject.Generic[3]   // T4
	log.Println(g0.Name, g0.OutObject.Type)
	log.Println(g1.Name, g1.OutObject.Type, g10.OutObject.Type, g100.Name, g1000.Name)
	log.Println(g2.Name, g2.OutObject.Type, g20.OutObject.Type, g21.OutObject.Type)
	log.Println(g3.Name, g3.OutObject.Type)
}
