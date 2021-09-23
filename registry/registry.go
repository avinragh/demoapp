package registry

import (
	"demoapp/context"
	"demoapp/evaluator"
	"sync"
)

var once sync.Once

const (
	ROOT_NODE_NAME = "ALARM_ROOT"
)

type rootNode struct{}

func (r *rootNode) GetName() string {
	return ROOT_NODE_NAME
}

var rootEntity *rootNode = &rootNode{}

type tree struct {
	root *node
}

type node struct {
	parent   *node
	name     string
	entity   evaluator.EvaluatorI
	children []*node
}

var singleton *tree

func Registry() *tree {
	if singleton == nil {
		once.Do(
			func() {
				singleton = &tree{
					root: &node{
						parent:   nil,
						name:     rootEntity.GetName(),
						entity:   rootEntity,
						children: make([]*node, 0),
					},
				}
			})
	}
	return singleton
}

func (t *tree) FindNode(name string) *node {
	queue := make([]*node, 0)
	queue = append(queue, t.root)
	for len(queue) > 0 {
		searchItem := queue[0]
		queue = queue[1:]
		if searchItem.name == name {
			return searchItem
		}
		if len(searchItem.children) > 0 {
			queue = append(queue, searchItem.children...)
		}
	}
	return nil
}

func (t *tree) AddNode(ctx *context.Context, aInEvaluator evaluator.EvaluatorI, aInParentEvaluator evaluator.EvaluatorI) {
	if existing := t.FindNode(aInEvaluator.GetName()); existing != nil {
		return
	}
	if aInParentEvaluator == nil {
		t.root.children = append(t.root.children, &node{
			parent:   t.root,
			entity:   aInEvaluator,
			name:     aInEvaluator.GetName(),
			children: make([]*node, 0),
		})
	} else if parentNode := t.FindNode(aInParentEvaluator.GetName()); parentNode != nil {
		parentNode.children = append(parentNode.children, &node{
			parent:   parentNode,
			entity:   aInEvaluator,
			name:     aInEvaluator.GetName(),
			children: make([]*node, 0),
		})
	}

}

func (t *tree) GetDescendents(item evaluator.EvaluatorI) []evaluator.EvaluatorI {
	var descendents []evaluator.EvaluatorI

	if itemNode := t.FindNode(item.GetName()); itemNode != nil {
		if len(itemNode.children) > 0 {
			for _, child := range itemNode.children {
				descendents = append(descendents, child.entity)
				if grandChildren := t.GetDescendents(child.entity); len(grandChildren) > 0 {
					descendents = append(descendents, grandChildren...)
				}
			}
		}
	}
	return descendents
}

func (t *tree) GetParent(item evaluator.EvaluatorI) evaluator.EvaluatorI {
	if item != nil {
		if itemNode := t.FindNode(item.GetName()); itemNode != nil {
			if itemNode.parent.entity.GetName() != ROOT_NODE_NAME {
				return itemNode.parent.entity
			}
			return nil
		}
	}
	return nil
}

func (t *tree) Length() int {
	return len(t.GetDescendents(rootEntity))
}

func (t *tree) Traverse(ctx *context.Context, evalFn func(curentEvaluator evaluator.EvaluatorI, descendents []evaluator.EvaluatorI) bool) {
	queue := make([]*node, 0)
	queue = append(queue, t.root.children...)
	for len(queue) > 0 {
		currentEvaluator := queue[0]
		queue = queue[1:]
		descendents := t.GetDescendents(currentEvaluator.entity)
		if evalFn(currentEvaluator.entity, descendents) {
			queue = append(queue, currentEvaluator.children...)
		}
	}
}
