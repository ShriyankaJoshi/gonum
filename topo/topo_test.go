// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
)

func TestIsPath(t *testing.T) {
	dg := concrete.NewDirectedGraph(0, math.Inf(1))
	if !IsPathIn(dg, nil) {
		t.Error("IsPath returns false on nil path")
	}
	p := []graph.Node{concrete.Node(0)}
	if IsPathIn(dg, p) {
		t.Error("IsPath returns true on nonexistant node")
	}
	dg.AddNode(p[0])
	if !IsPathIn(dg, p) {
		t.Error("IsPath returns false on single-length path with existing node")
	}
	p = append(p, concrete.Node(1))
	dg.AddNode(p[1])
	if IsPathIn(dg, p) {
		t.Error("IsPath returns true on bad path of length 2")
	}
	dg.SetEdge(concrete.Edge{F: p[0], T: p[1], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath returns false on correct path of length 2")
	}
	p[0], p[1] = p[1], p[0]
	if IsPathIn(dg, p) {
		t.Error("IsPath erroneously returns true for a reverse path")
	}
	p = []graph.Node{p[1], p[0], concrete.Node(2)}
	dg.SetEdge(concrete.Edge{F: p[1], T: p[2], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath does not find a correct path for path > 2 nodes")
	}
	ug := concrete.NewGraph(0, math.Inf(1))
	ug.SetEdge(concrete.Edge{F: p[1], T: p[0], W: 1})
	ug.SetEdge(concrete.Edge{F: p[1], T: p[2], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath does not correctly account for undirected behavior")
	}
}

var connectedComponentTests = []struct {
	g    []set
	want [][]int
}{
	{
		g: batageljZaversnikGraph,
		want: [][]int{
			{0},
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	},
}

func TestConnectedComponents(t *testing.T) {
	for i, test := range connectedComponentTests {
		g := concrete.NewGraph(0, math.Inf(1))

		for u, e := range test.g {
			if !g.Has(concrete.Node(u)) {
				g.AddNode(concrete.Node(u))
			}
			for v := range e {
				if !g.Has(concrete.Node(v)) {
					g.AddNode(concrete.Node(v))
				}
				g.SetEdge(concrete.Edge{F: concrete.Node(u), T: concrete.Node(v)})
			}
		}
		cc := ConnectedComponents(g)
		got := make([][]int, len(cc))
		for j, c := range cc {
			ids := make([]int, len(c))
			for k, n := range c {
				ids[k] = n.ID()
			}
			sort.Ints(ids)
			got[j] = ids
		}
		sort.Sort(internal.BySliceValues(got))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected connected components for test %d %T:\ngot: %v\nwant:%v", i, g, got, test.want)
		}
	}
}
