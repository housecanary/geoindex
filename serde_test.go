package geoindex

import (
	"encoding/binary"
	"github.com/tidwall/rbang"
	"io"
	"os"
	"testing"
)

func itemSaver(w io.Writer, value interface{}) (err error) {
	item := value.(uint32)
	if err = binary.Write(w, binary.BigEndian, item); err != nil {
		return
	}
	return
}

func itemLoader(r io.Reader, obuf []byte) (value interface{}, buf []byte, err error) {
	buf = obuf[:]
	var item uint32
	if err = binary.Read(r, binary.BigEndian, &item); err != nil {
		return
	}
	return item, buf,nil
}

func TestSaveLoadGeoIndex(t *testing.T) {
	N := 256
	boxes := randBoxes(N)

	tr := &rbang.RTree{}
	for i, box := range boxes {
		tr.Insert(box.min, box.max, uint32(i))
	}
	if tr.Len() != N {
		t.Fatalf("expected %d, got %d", N, tr.Len())
	}

	var f *os.File
	var err error
	fileName := "/tmp/tree_save"
	f, err = os.Create(fileName)
	if err != nil {
		t.Fatal("creating failed")
	}

	if err = tr.Save(f, itemSaver); err != nil {
		t.Fatal("saving failed")
	}
	if f.Close() != nil {
		t.Fatal("closing failed")
	}

	f, err = os.Open(fileName)
	if err != nil {
		t.Fatal("opening failed")
	}

	newTr := &rbang.RTree{}

	if err = newTr.Load(f, itemLoader); err != nil {
		t.Fatal("loading failed")
	}

	if f.Close() != nil {
		t.Fatal("closing failed")
	}

	if newTr.Len() != N {
		t.Fatalf("expected %d, got %d", N, newTr.Len())
	}

	var boxes1, boxes2 []uint32
	tr.Scan(func(min, max [2]float64, value interface{}) bool {
		boxes1 = append(boxes1, value.(uint32))
		return true
	})
	newTr.Scan(func(min, max [2]float64, value interface{}) bool {
		boxes2 = append(boxes2, value.(uint32))
		return true
	})
	for i := 0; i < len(boxes1); i++ {
		if boxes1[i] != boxes2[i] {
			t.Fatalf("expected '%v', got '%v'", boxes2[i], boxes1[i])
		}
	}
}
