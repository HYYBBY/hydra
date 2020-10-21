package container

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func TestNewContainer(t *testing.T) {
	l := NewContainer()
	if !reflect.DeepEqual(l, &Container{
		cache: cmap.New(8),
		vers:  newVers(),
	}) {
		t.Error("NewContainer() didn't return *Container")
	}
}

func TestContainer_GetOrCreate(t *testing.T) {
	
	c := NewContainer()
	type args struct {
		typ     string
		name    string
		creator func(conf *conf.RawConf) (interface{}, error)
	}
	creator := func(conf *conf.RawConf) (interface{}, error) {
		return nil, nil
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{typ: "db", name: "db", creator: creator}, want: "xxxx", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetOrCreate(tt.args.typ, tt.args.name, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("Container.GetOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Container.GetOrCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}