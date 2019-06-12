package ttlscheduler

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
)

func Test_scheduler_BindStream(t *testing.T) {
	type args struct {
		in0 *pb.Stream
	}
	tests := []struct {
		name     string
		nodes    map[int32]*pb.Node
		streams  []*pb.Stream
		bindings []*pb.Binding
		args     args
		want     *pb.Binding
		wantErr  bool
	}{
		{
			name:     "test 0",
			nodes:    map[int32]*pb.Node{},
			streams:  []*pb.Stream{},
			bindings: []*pb.Binding{},
			args:     args{&pb.Stream{Id: 1, MsgRate: 0}},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "test 1",
			nodes:    map[int32]*pb.Node{1: &pb.Node{Id: 1}},
			streams:  []*pb.Stream{},
			bindings: []*pb.Binding{},
			args:     args{&pb.Stream{Id: 1, MsgRate: 0}},
			want:     &pb.Binding{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
			streams:  []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}},
			bindings: []*pb.Binding{{Stream: &pb.Stream{Id: 1, MsgRate: 110}, Node: &pb.Node{Id: 1}}},
			args:     args{&pb.Stream{Id: 2, MsgRate: 0}},
			want:     &pb.Binding{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 2}},
			wantErr:  false,
		},
		{
			name:    "test 3",
			nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
			streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
			bindings: []*pb.Binding{
				{Stream: &pb.Stream{Id: 1, MsgRate: 110}, Node: &pb.Node{Id: 1}},
				{Stream: &pb.Stream{Id: 2, MsgRate: 50}, Node: &pb.Node{Id: 2}},
			},
			args:    args{&pb.Stream{Id: 3, MsgRate: 0}},
			want:    &pb.Binding{Stream: &pb.Stream{Id: 3, MsgRate: 0}, Node: &pb.Node{Id: 2}},
			wantErr: false,
		},
		{
			name:    "test 4",
			nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
			streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
			bindings: []*pb.Binding{
				{Stream: &pb.Stream{Id: 1, MsgRate: 50}, Node: &pb.Node{Id: 1}},
				{Stream: &pb.Stream{Id: 2, MsgRate: 110}, Node: &pb.Node{Id: 2}},
			},
			args:    args{&pb.Stream{Id: 3, MsgRate: 0}},
			want:    &pb.Binding{Stream: &pb.Stream{Id: 3, MsgRate: 0}, Node: &pb.Node{Id: 1}},
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
		nodes       map[int32]*pb.Node
		nNodes      []*pb.Node
		wantAdded   map[int32]*pb.Node
		wantDeleted map[int32]*pb.Node
	}{
		{
			name:        "test 1",
			nodes:       map[int32]*pb.Node{},
			nNodes:      []*pb.Node{&pb.Node{Id: 1}},
			wantAdded:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
			wantDeleted: map[int32]*pb.Node{},
		},
		{
			name:        "test 2",
			nodes:       map[int32]*pb.Node{1: &pb.Node{Id: 1}},
			nNodes:      []*pb.Node{},
			wantAdded:   map[int32]*pb.Node{},
			wantDeleted: map[int32]*pb.Node{1: &pb.Node{Id: 1}},
		},
		{
			name:        "test 3",
			nodes:       map[int32]*pb.Node{1: &pb.Node{Id: 1}, 4: &pb.Node{Id: 4}, 2: &pb.Node{Id: 2}},
			nNodes:      []*pb.Node{&pb.Node{Id: 1}, &pb.Node{Id: 3}, &pb.Node{Id: 2}},
			wantAdded:   map[int32]*pb.Node{3: &pb.Node{Id: 3}},
			wantDeleted: map[int32]*pb.Node{4: &pb.Node{Id: 4}},
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
		nodes    map[int32]*pb.Node
		streams  []*pb.Stream
		bindings []*pb.Binding
	}
	type args struct {
		nNodes []*pb.Node
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
				nodes:    map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams:  []*pb.Stream{{Id: 1, MsgRate: 0}},
				bindings: []*pb.Binding{{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}}},
			},
			args: args{nNodes: []*pb.Node{&pb.Node{Id: 1}, &pb.Node{Id: 2}}},
			want: fields{
				nodes:    map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
				streams:  []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}},
				bindings: []*pb.Binding{{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}}},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 2}},
				},
			},
			args: args{nNodes: []*pb.Node{&pb.Node{Id: 1}}},
			want: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 1}},
				},
			},
		},
		{
			name: "test 3",
			fields: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}, 3: &pb.Node{Id: 3}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 10}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 50}, Node: &pb.Node{Id: 2}},
					{Stream: &pb.Stream{Id: 3, MsgRate: 10}, Node: &pb.Node{Id: 3}},
				},
			},
			args: args{nNodes: []*pb.Node{&pb.Node{Id: 1}, &pb.Node{Id: 2}}},
			want: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 10}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 50}, Node: &pb.Node{Id: 2}},
					{Stream: &pb.Stream{Id: 3, MsgRate: 0}, Node: &pb.Node{Id: 1}},
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
		nodes    map[int32]*pb.Node
		streams  []*pb.Stream
		bindings []*pb.Binding
	}
	type args struct {
		nStreams []*pb.Stream
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
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 1}},
				},
			},
			args: args{nStreams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}}},
			want: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
				},
			},
		},
		{
			name: "test 2",
			fields: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 1}},
				},
			},
			args: args{nStreams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}}},
			want: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 3, MsgRate: 0}, Node: &pb.Node{Id: 1}},
				},
			},
		},
		{
			name: "test 3",
			fields: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 50}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 60}, Node: &pb.Node{Id: 2}},
				},
			},
			args: args{nStreams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}}},
			want: fields{
				nodes:   map[int32]*pb.Node{1: &pb.Node{Id: 1}, 2: &pb.Node{Id: 2}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}, &pb.Stream{Id: 3, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 50}, Node: &pb.Node{Id: 1}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 60}, Node: &pb.Node{Id: 2}},
					{Stream: &pb.Stream{Id: 3, MsgRate: 0}, Node: &pb.Node{Id: 1}},
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
		nodes map[int32]*pb.Node
	}
	tests := []struct {
		name  string
		nodes map[int32]*pb.Node
		want  map[int32]*pb.Node
	}{
		{
			name:  "test 1",
			nodes: map[int32]*pb.Node{1: &pb.Node{Id: 1}},
			want:  map[int32]*pb.Node{1: &pb.Node{Id: 1}},
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

func (tnf *testNodeFetcher) GetNodes() ([]*pb.Node, error) {
	return []*pb.Node{&pb.Node{Id: 0}, &pb.Node{Id: 1}}, nil
}

type testStreamFetcher struct{}

func (tsf *testStreamFetcher) GetStreams() ([]*pb.Stream, error) {
	return []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}}, nil
}

func Test_basicScheduler_Observe(t *testing.T) {
	t.Log("!! Note this test will fail sometimes due to randomized order of node map !!")

	type fields struct {
		nodes    map[int32]*pb.Node
		streams  []*pb.Stream
		bindings []*pb.Binding
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
				nodes:    map[int32]*pb.Node{},
				streams:  []*pb.Stream{},
				bindings: []*pb.Binding{},
				nf:       &testNodeFetcher{},
				sf:       &testStreamFetcher{},
				once:     true,
			},
			want: fields{
				nodes:   map[int32]*pb.Node{0: &pb.Node{Id: 0}, 1: &pb.Node{Id: 1}},
				streams: []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}, &pb.Stream{Id: 2, MsgRate: 0}},
				bindings: []*pb.Binding{
					{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 0}},
					{Stream: &pb.Stream{Id: 2, MsgRate: 0}, Node: &pb.Node{Id: 1}},
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
		nodes    map[int32]*pb.Node
		streams  []*pb.Stream
		bindings []*pb.Binding
		req      *pb.NodeBindingRequest
		want     *pb.NodeBindingResponse
		wantErr  bool
	}{
		{
			name:     "test 1",
			nodes:    map[int32]*pb.Node{0: &pb.Node{Id: 0}},
			streams:  []*pb.Stream{},
			bindings: []*pb.Binding{},
			req:      &pb.NodeBindingRequest{Node: &pb.Node{Id: 0}},
			want:     &pb.NodeBindingResponse{Bindings: []*pb.Binding{}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    map[int32]*pb.Node{0: &pb.Node{Id: 0}},
			streams:  []*pb.Stream{&pb.Stream{Id: 1, MsgRate: 0}},
			bindings: []*pb.Binding{&pb.Binding{Stream: &pb.Stream{Id: 1, MsgRate: 0}, Node: &pb.Node{Id: 0}}},
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
