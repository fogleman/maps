package maps

import (
	"io"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
)

type PBF struct {
	Nodes     map[int64]*osmpbf.Node
	Ways      map[int64]*osmpbf.Way
	Relations map[int64]*osmpbf.Relation
}

func LoadPBF(path string) (*PBF, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := osmpbf.NewDecoder(file)
	err = decoder.Start(runtime.GOMAXPROCS(0))
	if err != nil {
		return nil, err
	}

	nodes := make(map[int64]*osmpbf.Node)
	ways := make(map[int64]*osmpbf.Way)
	relations := make(map[int64]*osmpbf.Relation)

	for {
		if v, err := decoder.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				nodes[v.ID] = v
			case *osmpbf.Way:
				ways[v.ID] = v
			case *osmpbf.Relation:
				relations[v.ID] = v
			}
		}
	}

	pbf := PBF{nodes, ways, relations}
	return &pbf, nil
}
