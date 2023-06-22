package helper

import (
	"context"
	"github.com/whatvn/denny"
	"github.com/whatvn/denny/log"
)

const ActorCtxKey = "actor"

func getActorFromContext(context context.Context) (string, bool) {
	if ctx, ok := context.(*denny.Context); ok {
		iActor, ok := ctx.Get(ActorCtxKey)
		if !ok {
			return "", false
		}
		return iActor.(string), true
	}

	iActor := context.Value(ActorCtxKey)
	if iActor == nil {
		return "", false
	}
	return iActor.(string), true
}

func GetUserAndLogger(ctx context.Context) (string, *log.Log) {
	actor, _ := getActorFromContext(ctx)
	logger := denny.GetLogger(ctx).WithField("actor", actor)
	return actor, logger
}
