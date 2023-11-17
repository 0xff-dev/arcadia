package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"strings"

	"github.com/kubeagi/arcadia/graphql-server/go-server/graph/model"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/auth"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/client"
	md "github.com/kubeagi/arcadia/graphql-server/go-server/pkg/model"
)

// CreateModel is the resolver for the createModel field.
func (r *mutationResolver) CreateModel(ctx context.Context, input model.CreateModelInput) (*model.Model, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}

	displayname, description, filed, modeltype := "", "", "", ""

	if input.DisplayName != "" {
		displayname = input.DisplayName
	}
	if input.Description != nil {
		description = *input.Description
	}
	if input.Field != "" {
		filed = input.Field
	}
	if input.Modeltype != "" {
		modeltype = input.Modeltype
	}
	return md.CreateModel(ctx, c, input.Name, input.Namespace, displayname, description, filed, modeltype)
}

// UpdateModel is the resolver for the updateModel field.
func (r *mutationResolver) UpdateModel(ctx context.Context, input *model.UpdateModelInput) (*model.Model, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name, displayname := "", ""
	if input.DisplayName != "" {
		displayname = input.DisplayName
	}
	if input.Name != "" {
		name = input.Name

	}
	return md.UpdateModel(ctx, c, name, input.Namespace, displayname)
}

// DeleteModel is the resolver for the deleteModel field.
func (r *mutationResolver) DeleteModel(ctx context.Context, input *model.DeleteModelInput) (*string, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name := ""
	labelSelector, fieldSelector := "", ""
	if input.Name != nil {
		name = *input.Name
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	return md.DeleteModel(ctx, c, name, input.Namespace, labelSelector, fieldSelector)
}

// GetModel is the resolver for the getModel field.
func (r *queryResolver) GetModel(ctx context.Context, name string, namespace string) (*model.Model, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	return md.ReadModel(ctx, c, name, namespace)
}

// ListModels is the resolver for the listModels field.
func (r *queryResolver) ListModels(ctx context.Context, input model.ListModelInput) (*model.PaginatedModel, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name, displayName, labelSelector, fieldSelector := "", "", "", ""
	page, pageSize := 1, 10
	if input.Name != nil {
		name = *input.Name
	}
	if input.DisplayName != nil {
		displayName = *input.DisplayName
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	if input.Page != nil && *input.Page > 0 {
		page = *input.Page
	}
	if input.PageSize != nil && *input.PageSize > 0 {
		pageSize = *input.PageSize
	}
	result, err := md.ListModels(ctx, c, input.Namespace, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}
	var filteredResult []*model.Model
	for idx, u := range result {
		if (name == "" || strings.Contains(u.Name, name)) && (displayName == "" || strings.Contains(u.DisplayName, displayName)) {
			filteredResult = append(filteredResult, result[idx])
		}
	}

	totalCount := len(filteredResult)
	end := page * pageSize
	if end > totalCount {
		end = totalCount
	}
	return &model.PaginatedModel{
		TotalCount:  totalCount,
		HasNextPage: end < totalCount,
		Nodes:       filteredResult[(page-1)*pageSize : end],
	}, nil
}