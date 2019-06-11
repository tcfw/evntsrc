package ttlscheduler

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

func Test_scheduler_BindStream(t *testing.T) {
	type args struct {
		in0 *stream
	}
	tests := []struct {
		name     string
		nodes    map[int]*node
		streams  []*stream
		bindings []*binding
		args     args
		want     *binding
		wantErr  bool
	}{
		{
			name:     "test 0",
			nodes:    map[int]*node{},
			streams:  []*stream{},
			bindings: []*binding{},
			args:     args{&stream{1, 0}},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "test 1",
			nodes:    map[int]*node{1: {1}},
			streams:  []*stream{},
			bindings: []*binding{},
			args:     args{&stream{1, 0}},
			want:     &binding{Stream: &stream{1, 0}, Node: &node{1}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    map[int]*node{1: {1}, 2: {2}},
			streams:  []*stream{{1, 0}},
			bindings: []*binding{{Stream: &stream{1, 110}, Node: &node{1}}},
			args:     args{&stream{2, 0}},
			want:     &binding{Stream: &stream{2, 0}, Node: &node{2}},
			wantErr:  false,
		},
		{
			name:    "test 3",
			nodes:   map[int]*node{1: {1}, 2: {2}},
			streams: []*stream{{1, 0}, {2, 0}},
			bindings: []*binding{
				{Stream: &stream{1, 110}, Node: &node{1}},
				{Stream: &stream{2, 50}, Node: &node{2}},
			},
			args:    args{&stream{3, 0}},
			want:    &binding{Stream: &stream{3, 0}, Node: &node{2}},
			wantErr: false,
		},
		{
			name:    "test 4",
			nodes:   map[int]*node{1: {1}, 2: {2}},
			streams: []*stream{{1, 0}, {2, 0}},
			bindings: []*binding{
				{Stream: &stream{1, 50}, Node: &node{1}},
				{Stream: &stream{2, 110}, Node: &node{2}},
			},
			args:    args{&stream{3, 0}},
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
		streams  []*stream
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
				streams:  []*stream{{1, 0}},
				bindings: []*binding{{Stream: &stream{1, 0}, Node: &node{1}}},
			},
			args: args{nNodes: []*node{{1}, {2}}},
			want: fields{
				nodes:    map[int]*node{1: {1}, 2: {2}},
				streams:  []*stream{{1, 0}},
				bindings: []*binding{{Stream: &stream{1, 0}, Node: &node{1}}},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []*stream{{1, 0}, {2, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{2}},
				},
			},
			args: args{nNodes: []*node{{1}}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []*stream{{1, 0}, {2, 0}},
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
				streams: []*stream{{1, 0}, {2, 0}, {3, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 10}, Node: &node{1}},
					{Stream: &stream{2, 50}, Node: &node{2}},
					{Stream: &stream{3, 10}, Node: &node{3}},
				},
			},
			args: args{nNodes: []*node{{1}, {2}}},
			want: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []*stream{{1, 0}, {2, 0}, {3, 0}},
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
		streams  []*stream
		bindings []*binding
	}
	type args struct {
		nStreams []*stream
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
				streams: []*stream{{1, 0}, {2, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
			},
			args: args{nStreams: []*stream{{1, 0}}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []*stream{{1, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
				},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []*stream{{1, 0}, {2, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{1}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
			},
			args: args{nStreams: []*stream{{1, 0}, {2, 0}, {3, 0}}},
			want: fields{
				nodes:   map[int]*node{1: {1}},
				streams: []*stream{{1, 0}, {2, 0}, {3, 0}},
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
				streams: []*stream{{1, 0}, {2, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 50}, Node: &node{1}},
					{Stream: &stream{2, 60}, Node: &node{2}},
				},
			},
			args: args{nStreams: []*stream{{1, 0}, {2, 0}, {3, 0}}},
			want: fields{
				nodes:   map[int]*node{1: {1}, 2: {2}},
				streams: []*stream{{1, 0}, {2, 0}, {3, 0}},
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

func Test_basicScheduler_GetNodes(t *testing.T) {
	type fields struct {
		nodes map[int]*node
	}
	tests := []struct {
		name  string
		nodes map[int]*node
		want  map[int]*node
	}{
		{
			name:  "test 1",
			nodes: map[int]*node{1: {1}},
			want:  map[int]*node{1: {1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes: tt.nodes,
			}
			if got := s.GetNodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("basicScheduler.GetNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testNodeFetcher struct{}

func (tnf *testNodeFetcher) GetNodes() ([]*node, error) {
	return []*node{{0}, {1}}, nil
}

type testStreamFetcher struct{}

func (tsf *testStreamFetcher) GetStreams() ([]*stream, error) {
	return []*stream{{1, 0}, {2, 0}}, nil
}

func Test_basicScheduler_Observe(t *testing.T) {
	t.Log("!! Note this test will fail sometimes due to randomized order of node map !!")

	type fields struct {
		nodes    map[int]*node
		streams  []*stream
		bindings []*binding
		nf       NodeFetcher
		sf       StreamFetcher
		once     bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    fields
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				nodes:    map[int]*node{},
				streams:  []*stream{},
				bindings: []*binding{},
				nf:       &testNodeFetcher{},
				sf:       &testStreamFetcher{},
				once:     true,
			},
			want: fields{
				nodes:   map[int]*node{0: {0}, 1: {1}},
				streams: []*stream{{1, 0}, {2, 0}},
				bindings: []*binding{
					{Stream: &stream{1, 0}, Node: &node{0}},
					{Stream: &stream{2, 0}, Node: &node{1}},
				},
				nf:   &testNodeFetcher{},
				sf:   &testStreamFetcher{},
				once: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.fields.nodes,
				streams:  tt.fields.streams,
				bindings: tt.fields.bindings,
				nf:       tt.fields.nf,
				sf:       tt.fields.sf,
				once:     tt.fields.once,
			}
			err := s.Observe()
			ttComp := fields{
				nodes:    s.nodes,
				streams:  s.streams,
				bindings: s.bindings,
				nf:       s.nf,
				sf:       s.sf,
				once:     s.once,
			}
			if !reflect.DeepEqual(ttComp, tt.want) {
				t.Errorf("basicScheduler.Observe() = %v, want %v", ttComp, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("basicScheduler.Observe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_basicScheduler_NodeBindings(t *testing.T) {
	tests := []struct {
		name     string
		nodes    map[int]*node
		streams  []*stream
		bindings []*binding
		req      *pb.NodeBindingRequest
		want     *pb.NodeBindingResponse
		wantErr  bool
	}{
		{
			name:     "test 1",
			nodes:    map[int]*node{0: {0}},
			streams:  []*stream{},
			bindings: []*binding{},
			req:      &pb.NodeBindingRequest{Node: &pb.Node{Id: 0}},
			want:     &pb.NodeBindingResponse{Bindings: []*pb.Binding{}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    map[int]*node{0: {0}},
			streams:  []*stream{&stream{1, 0}},
			bindings: []*binding{&binding{Stream: &stream{1, 0}, Node: &node{0}}},
			req:      &pb.NodeBindingRequest{Node: &pb.Node{Id: 0}},
			want:     &pb.NodeBindingResponse{Bindings: []*pb.Binding{&pb.Binding{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 0}}}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &basicScheduler{
				nodes:    tt.nodes,
				streams:  tt.streams,
				bindings: tt.bindings,
			}
			got, err := s.NodeBindings(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("basicScheduler.NodeBindings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("basicScheduler.NodeBindings() = %v, want %v", got, tt.want)
			}
		})
	}
}
