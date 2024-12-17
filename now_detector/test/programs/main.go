package programs

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/spanner"
)

func main() {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	err = insert(ctx, client, true)
	if err != nil {
		log.Fatal(err)
	}
}

func insert(ctx context.Context, client *spanner.Client, isNow bool) error {
	var now time.Time
	if isNow {
		now = time.Now()
	} else {
		now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	_, err := client.Apply(ctx, []*spanner.Mutation{
		spanner.Insert(
			"Users",
			[]string{"name", "created_at"},
			[]interface{}{"Alice", now},
		),
	})

	return err
}
