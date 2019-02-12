package main

import (
	"context"
	"fmt"
	"log"

	"flag"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers               = []string{"localhost:9092"}
	topic     goka.Stream = "example-stream"
	group     goka.Group  = "example-group"
	emitter               = flag.Bool("emitter", false, "run emitter")
	processor             = flag.Bool("processor", false, "run processor")
)

// emits a single message and leave
func runEmitter() {
	emitter, err := goka.NewEmitter(brokers, topic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitter.Finish()

	err = emitter.EmitSync("hi", "hello")
	if err != nil {
		log.Fatalf("error emitting message: %v", err)
	}
	fmt.Println("message emitted")
}

// process messages until ctrl-c is pressed
func runProcessor() {
	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(ctx goka.Context, msg interface{}) {

		// during the second run, this should break (as value should already be in context)
		if val := ctx.Value(); val != nil {
			panic(fmt.Sprintf("dealing with a value already in context %v", ctx.Value()))
		}

		// store received value in context (first run)
		ctx.SetValue(msg.(string))
		log.Printf("stored to ctx key = %s, msg = %v", ctx.Key(), msg)
	}

	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.String), cb),
		goka.Persist(new(codec.String)),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
	if err = p.Run(context.Background()); err != nil {
		log.Fatalf("error running processor: %v", err)
	}
}

func main() {
	flag.Parse()
	if *processor {
		runProcessor() // press ctrl-c to stop
	}
	if *emitter {
		runEmitter() // emits one message and stops
	}
}
