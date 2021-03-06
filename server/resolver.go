package server

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/adyatlov/bunxp/xp"
	"github.com/graph-gophers/graphql-go"
)

type resolver struct {
	explorer *xp.Explorer
}

func (r *resolver) Root() (*objectResolver, error) {
	root, err := r.explorer.Root()
	if err != nil {
		return nil, err
	}
	return &objectResolver{object: root}, nil
}

func (r *resolver) Object(args struct {
	ObjectId string
	Type     string
}) (*objectResolver, error) {
	t := xp.ObjectTypeName(args.Type)
	id := xp.ObjectId(args.ObjectId)
	object, err := r.explorer.Object(t, id)
	return &objectResolver{object: object}, err
}

func (r *resolver) ObjectTypes() []*objectTypeResolver {
	objTypeMap := xp.GetObjectTypes()
	objTypes := make([]*objectTypeResolver, 0, len(objTypeMap))
	for _, t := range objTypeMap {
		objTypes = append(objTypes, &objectTypeResolver{t: t})
	}
	return objTypes
}

func (r *resolver) MetricTypes() *[]*metricTypeResolver {
	mTypeMap := xp.GetMetricTypes()
	mTypes := make([]*metricTypeResolver, 0, len(mTypeMap))
	for _, t := range mTypeMap {
		mTypes = append(mTypes, &metricTypeResolver{t: t})
	}
	return &mTypes
}

type objectResolver struct {
	object xp.Object
}

func (r *objectResolver) ID() graphql.ID {
	idSting := []byte(string(r.object.Type()) + ":::" + string(r.object.Id()))
	return graphql.ID(base64.StdEncoding.EncodeToString(idSting))
}

func (r *objectResolver) ObjectId() string {
	return string(r.object.Id())
}

func (r *objectResolver) Type() string {
	return string(r.object.Type())
}

func (r *objectResolver) Name() string {
	return string(r.object.Name())
}

func (r *objectResolver) Metrics(args struct {
	Names []string
}) (*[]*metricResolver, error) {
	metricTypeNames := make([]xp.MetricTypeName, 0, len(args.Names))
	for _, m := range args.Names {
		metricTypeNames = append(metricTypeNames, xp.MetricTypeName(m))
	}
	metrics, err := r.object.Metrics(metricTypeNames...)
	if err != nil {
		return nil, err
	}
	metricResolvers := make([]*metricResolver, 0, len(metrics))
	for _, metric := range metrics {
		metricResolvers = append(metricResolvers,
			&metricResolver{metric: metric})
	}
	return &metricResolvers, nil
}

func (r *objectResolver) Children() (*[]*objectGroupResolver, error) {
	children, err := r.object.Children()
	if err != nil {
		return nil, err
	}
	objectGroupResolvers := make([]*objectGroupResolver, 0, len(children))
	for _, group := range children {
		objectGroupResolvers = append(objectGroupResolvers, &objectGroupResolver{group: &group})
	}
	return &objectGroupResolvers, nil
}

type objectGroupResolver struct {
	group *xp.ObjectGroup
}

func (r *objectGroupResolver) Type() string {
	return string(r.group.Type)
}

func (r *objectGroupResolver) Objects() *[]*objectResolver {
	objectResolvers := make([]*objectResolver, 0, len(r.group.Objects))
	for _, object := range r.group.Objects {
		objectResolvers = append(objectResolvers, &objectResolver{object: object})
	}
	return &objectResolvers
}

type metricResolver struct {
	metric *xp.Metric
}

func (r *metricResolver) Type() string {
	return string(r.metric.Type)
}

func (r *metricResolver) Value() string {
	return fmt.Sprintf("%v", r.metric.Value)
}

type objectTypeResolver struct {
	t xp.ObjectType
}

func (r *objectTypeResolver) Name() string {
	return string(r.t.Name)
}

func (r *objectTypeResolver) DisplayName() string {
	return r.t.DisplayName
}

func (r *objectTypeResolver) PluralDisplayName() string {
	return r.t.PluralDisplayName
}

func (r *objectTypeResolver) Description() string {
	return r.t.Description
}

func (r *objectTypeResolver) Metrics() *[]string {
	metrics := make([]string, 0, len(r.t.Metrics))
	for _, m := range r.t.Metrics {
		metrics = append(metrics, string(m))
	}
	return &metrics
}

func (r *objectTypeResolver) DefaultMetrics() *[]string {
	metrics := make([]string, 0, len(r.t.DefaultMetrics))
	for _, m := range r.t.DefaultMetrics {
		metrics = append(metrics, string(m))
	}
	return &metrics
}

type metricTypeResolver struct {
	t xp.MetricType
}

func (r *metricTypeResolver) Name() string {
	return string(r.t.Name)
}

func (r *metricTypeResolver) ObjectTypeName() string {
	return string(r.t.ObjectTypeName)
}

func (r *metricTypeResolver) ValueType() string {
	return strings.ToUpper(string(r.t.ValueType))
}

func (r *metricTypeResolver) DisplayName() string {
	return r.t.DisplayName
}

func (r *metricTypeResolver) Description() string {
	return r.t.Description
}
