package main

import (
	"log"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// cache = redis.Redis(host='redis', port=6379)
/**

def get_hit_count():
    retries = 5
    while True:
        try:
            return cache.incr('hits')
        except redis.exceptions.ConnectionError as exc:
            if retries == 0:
                raise exc
            retries -= 1
            time.sleep(0.5)


@app.route('/')
def hello():
    count = get_hit_count()
    return 'Hello World! I have been seen {} times.\n'.format(count)

**/

func main() {
	driver, err := neo4j.NewDriver("bolt://neo4j:7687", neo4j.BasicAuth("neo4j", "test", "")) /*func(c *neo4j.Config) {
		// https://neo4j.com/developer/docker-run-neo4j/
		// By default, the docker image does not have certificates installed.
		// This means that you will need to disable encryption when connecting with a driver.
		c.Encrypted = false
	}*/

	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		defer session.Close()

		greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
				map[string]interface{}{"message": "hello, world"})

			if err != nil {
				return nil, err
			}

			if result.Next() {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})

		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.JSON(iris.Map{
				"code": http.StatusOK,
				"data": greeting.(string),
			})
		}

		// for result.Next() {
		// 	record = result.Record();
		// 	if value, ok := record.Get('field_name'); ok {
		// 		// a value with alias field_name was found
		// 		// process value
		// 	}
		// }
	})

	err = app.Run(
		// Start the web server at localhost:5000
		iris.Addr(":5000"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)

	if err != nil {
		log.Fatal(err)
	}
}
