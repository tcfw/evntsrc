package ttlscheduler

import (
	"reflect"
	"testing"
)

func Test_scheduler_NodeBindings(t *testing.T) {
	type args struct {
		node node
	}
	tests := []struct {
		name     string
		nodes    map[int]*node
		streams  []int32
		bindings []*binding
		args     args
		want     []*binding
		wantErr  bool
	}{
		{
			name:     "test 1",
			nodes:    map[int]*node{0: {0}},
			streams:  []int32{1},
			bindings: []*binding{{Stream: &stream{ID: 1}, Node: &node{0}}},
			args:     args{node{0}},
			want:     []*binding{{Stream: &stream{ID: 1}, Node: &node{0}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.nodes,
				streams:  tt.streams,
				bindings: tt.bindings,
			}
			got, err := s.NodeBindings(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("scheduler.NodeBindings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scheduler.NodeBindings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scheduler_BindStream(t *testing.T) {
	type args struct {
		in0 int32
	}
	tests := []struct {
		name     string
		nodes    map[int]*node
		streams  []int32
		bindings []*binding
		args     args
		want     *binding
		wantErr  bool
	}{
		{
			name:     "test 0",
			nodes:    map[int]*node{},
			streams:  []int32{},
			bindings: []*binding{},
			args:     args{1},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "test 1",
			nodes:    map[int]*node{1: {1}},
			streams:  []int32{},
			bindings: []*binding{},
			args:     args{1},
			want:     &binding{Stream: &stream{1, 0}, Node: &node{1}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    map[int]*node{1: {1}, 2: {2}},
			streams:  []int32{1},
			bindings: []*binding{{Stream: &stream{1, 110}, Node: &node{1}}},
			args:     args{2},
			want:     &binding{Stream: &stream{2, 0}, Node: &node{2}},
			wantErr:  false,
		},
		{
			name:    "test 3",
			nodes:   map[int]*node{1: {1}, 2: {2}},
			streams: []int32{1, 2},
			bindings: []*binding{
				{Stream: &stream{1, 110}, Node: &node{1}},
				{Stream: &stream{2, 50}, Node: &node{2}},
			},
			args:    args{3},
			want:    &binding{Stream: &stream{3, 0}, Node: &node{2}},
			wantErr: false,
		},
		{
			name:    "test 4",
			nodes:   map[int]*node{1: {1}, 2: {2}},
			streams: []int32{1, 2},
			bindings: []*binding{
				{Stream: &stream{1, 50}, Node: &node{1}},
				{Stream: &stream{2, 110}, Node: &node{2}},
			},
			args:    args{3},
			want:    &binding{Stream: &stream{3, 0}, Node: &node{1}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.nodes,
				streams:  tt.streams,
				bindings: tt.bindings,
			}
			got, err := s.BindStream(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("scheduler.BindStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scheduler.BindStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicScheduler_nodeDiff(t *testing.T) {
	tests := []struct {
		name        string
		nodes       map[int]*node
		nNodes      []*node
		wantAdded   map[int]*node
		wantDeleted map[int]*node
	}{
		{
			name:        "test 1",
			nodes:       map[int]*node{},
			nNodes:      []*node{{1}},
			wantAdded:   map[int]*node{1: {1}},
			wantDeleted: map[int]*node{},
		},
		{
			name:        "test 2",
			nodes:       map[int]*node{1: {1}},
			nNodes:      []*node{},
			wantAdded:   map[int]*node{},
			wantDeleted: map[int]*node{1: {1}},
		},
		{
			name:        "test 3",
			nodes:       map[int]*node{1: {1}, 4: {4}, 2: {2}},
			nNodes:      []*node{{1}, {3}, {2}},
			wantAdded:   map[int]*node{3: {3}},
			wantDeleted: map[int]*node{4: {4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes: tt.nodes,
			}
			got, got1 := s.nodeDiff(tt.nNodes)
			if !reflect.DeepEqual(got, tt.wantAdded) {
				t.Errorf("basicScheduler.nodeDiff() got = %v, want %v", got, tt.wantAdded)
			}
			if !reflect.DeepEqual(got1, tt.wantDeleted) {
				t.Errorf("basicScheduler.nodeDiff() got1 = %v, want %v", got1, tt.wantDeleted)
			}
		})
	}
}

func Test_basicScheduler_observeNodes(t *testing.T) {
	type fields struct {
		nodes    map[int]*node
		streams  []int32
		bindings []*binding
	}
	type args struct {
		nNodes []*node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "test 1",
			fields: fields{
				nodes:    map[int]*node{1: {1}},
				streams:  []int32{1},
				bindings: []*binding{{Stream: &stream{1, 0}, Node: &node{1}}},
			},
			args: args{nNodes: []*node{{1}, {2}}},
			want: fields{
				nodes:    map[int]*node{1: {1}, 2: {2}},
				streams:  []int32{1},
				bindings: []*binding{{Stream: &stream{1, 0}, Node: &node{1}}},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []int32{1, 2},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{2}},
				},
			},
			args: args{nNodes: []*node{{1}}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []int32{1, 2},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
			},
		},
		{
			name: "test 3",
			fields: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}, 3: {3}},
				streams: []int32{1, 2, 3},
				bindings: []*binding{
					{Stream: &stream{1, 10}, Node: &node{1}},
					{Stream: &stream{2, 50}, Node: &node{2}},
					{Stream: &stream{3, 10}, Node: &node{3}},
				},
			},
			args: args{nNodes: []*node{{1}, {2}}},
			want: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []int32{1, 2, 3},
				bindings: []*binding{
					{Stream: &stream{1, 10}, Node: &node{1}},
					{Stream: &stream{2, 50}, Node: &node{2}},
					{Stream: &stream{3, 0}, Node: &node{1}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.fields.nodes,
				streams:  tt.fields.streams,
				bindings: tt.fields.bindings,
			}
			s.observeNodes(tt.args.nNodes)
			ttComp := fields{
				nodes:    s.nodes,
				streams:  s.streams,
				bindings: s.bindings,
			}
			if !reflect.DeepEqual(ttComp, tt.want) {
				t.Errorf("basicScheduler.observeNodes() s = %v, want %v", ttComp, tt.want)
			}
		})
	}
}

func Test_basicScheduler_observeStreams(t *testing.T) {
	type fields struct {
		nodes    map[int]*node
		streams  []int32
		bindings []*binding
	}
	type args struct {
		nStreams []int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "test 1",
			fields: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []int32{1, 2},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
			},
			args: args{nStreams: []int32{1}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []int32{1},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
				},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []int32{1, 2},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
			},
			args: args{nStreams: []int32{1, 2, 3}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []int32{1, 2, 3},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
					{Stream: &stream{3, 0}, Node: &node{1}},
				},
			},
		},
		{
			name: "test 3",
			fields: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []int32{1, 2},
				bindings: []*binding{
					{Stream: &stream{1, 50}, Node: &node{1}},
					{Stream: &stream{2, 60}, Node: &node{2}},
				},
			},
			args: args{nStreams: []int32{1, 2, 3}},
			want: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []int32{1, 2, 3},
				bindings: []*binding{
					{Stream: &stream{1, 50}, Node: &node{1}},
					{Stream: &stream{2, 60}, Node: &node{2}},
					{Stream: &stream{3, 0}, Node: &node{1}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.fields.nodes,
				streams:  tt.fields.streams,
				bindings: tt.fields.bindings,
			}
			s.observeStreams(tt.args.nStreams)
			ttComp := fields{
				nodes:    s.nodes,
				streams:  s.streams,
				bindings: s.bindings,
			}
			if !reflect.DeepEqual(ttComp, tt.want) {
				t.Errorf("basicScheduler.observeStreams() s = %v, want %v", ttComp, tt.want)
			}
		})
	}
}
