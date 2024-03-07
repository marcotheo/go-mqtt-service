package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wind-c/comqtt/v2/mqtt"
	"github.com/wind-c/comqtt/v2/mqtt/hooks/auth"
	"github.com/wind-c/comqtt/v2/mqtt/listeners"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// An example of configuring various server options...
	options := &mqtt.Options{
		// InflightTTL: 60 * 15, // Set an example custom 15-min TTL for inflight messages
	}

	server := mqtt.New(options)

	// For security reasons, the default implementation disallows all connections.
	// If you want to allow all connections, you must specifically allow it.
	// err := server.AddHook(new(auth.Hook), &auth.Options{
	// 	Ledger: &auth.Ledger{
	// 		Auth: auth.AuthRules{ // Auth disallows all by default
	// 			{Username: "peach", Password: "password1", Allow: true},
	// 			{Username: "melon", Password: "password2", Allow: true},
	// 			{Remote: "127.0.0.1:*", Allow: true},
	// 			{Remote: "localhost:*", Allow: true},
	// 		},
	// 		ACL: auth.ACLRules{ // ACL allows all by default
	// 			{Remote: "127.0.0.1:*"}, // local superuser allow all
	// 			{
	// 				// user melon can read and write to their own topic
	// 				Username: "melon", Filters: auth.Filters{
	// 					"melon/#":   auth.ReadWrite,
	// 					"updates/#": auth.WriteOnly, // can write to updates, but can't read updates from others
	// 				},
	// 			},
	// 			{
	// 				// Otherwise, no clients have publishing permissions
	// 				Filters: auth.Filters{
	// 					"#":         auth.ReadOnly,
	// 					"updates/#": auth.Deny,
	// 				},
	// 			},
	// 		},
	// 	},
	// })

	err := server.AddHook(new(auth.AllowHook), nil)

	if err != nil {
		log.Fatal(err)
	}

	tcp := listeners.NewWebsocket("ws1", ":1882", nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	server.Log.Info("main.go finished")
}
