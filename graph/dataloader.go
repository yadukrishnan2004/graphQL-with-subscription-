package graph

import (
	"context"
	"gql/graph/model"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"
	"gorm.io/gorm"
)

type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model.User]
}

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

func DataLoaderMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := &Loaders{
			UserLoader: dataloadgen.NewLoader(func(ctx context.Context, userIds []string) ([]*model.User, []error) {
				var users []*model.User
				if err := db.Where("id IN ?", userIds).Find(&users).Error; err != nil {
					return nil, []error{err}
				}

				userMap := make(map[string]*model.User, len(users))
				for _, u := range users {
					userMap[u.ID] = u
				}

				result := make([]*model.User, len(userIds))
				for i, id := range userIds {
					result[i] = userMap[id]
				}
				return result, nil
			}, dataloadgen.WithWait(time.Millisecond*10)),
		}

		ctx := context.WithValue(r.Context(), ctxKey, loaders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func For(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(ctxKey).(*Loaders)
	if !ok {
		return nil
	}
	return loaders
}