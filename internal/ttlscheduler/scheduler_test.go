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
		nodes    []*node
		streams  []int32
		bindings []*binding
		args     args
		want     []*binding
		wantErr  bool
	}{
		{
			name:     "test 1",
			nodes:    []*node{&node{0}},
			streams:  []int32{1},
			bindings: []*binding{&binding{Stream: &stream{ID: 1}, Node: &node{0}}},
			args:     args{node{0}},
			want:     []*binding{&binding{Stream: &stream{ID: 1}, Node: &node{0}}},
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
		nodes    []*node
		streams  []int32
		bindings []*binding
		args     args
		want     *binding
		wantErr  bool
	}{
		{
			name:     "test 0",
			nodes:    []*node{},
			streams:  []int32{},
			bindings: []*binding{},
			args:     args{1},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "test 1",
			nodes:    []*node{&node{1}},
			streams:  []int32{},
			bindings: []*binding{},
			args:     args{1},
			want:     &binding{Stream: &stream{1, 0}, Node: &node{1}},
			wantErr:  false,
		},
		{
			name:     "test 2",
			nodes:    []*node{&node{1}, &node{2}},
			streams:  []int32{1},
			bindings: []*binding{&binding{Stream: &stream{1, 110}, Node: &node{1}}},
			args:     args{2},
			want:     &binding{Stream: &stream{2, 0}, Node: &node{2}},
			wantErr:  false,
		},
		{
			name:    "test 3",
			nodes:   []*node{&node{1}, &node{2}},
			streams: []int32{1, 2},
			bindings: []*binding{
				&binding{Stream: &stream{1, 110}, Node: &node{1}},
				&binding{Stream: &stream{2, 50}, Node: &node{2}},
			},
			args:    args{3},
			want:    &binding{Stream: &stream{3, 0}, Node: &node{2}},
			wantErr: false,
		},
		{
			name:    "test 4",
			nodes:   []*node{&node{1}, &node{2}},
			streams: []int32{1, 2},
			bindings: []*binding{
				&binding{Stream: &stream{1, 50}, Node: &node{1}},
				&binding{Stream: &stream{2, 110}, Node: &node{2}},
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
