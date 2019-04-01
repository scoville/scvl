package redis

import "testing"

const testRedisDB = 999

func TestRedis(t *testing.T) {
	var err error
	client, err := NewClient()
	if err != nil {
		t.Fatal("cannot connect to redis server")
	}
	defer client.Close()
	client.SetDatabase(testRedisDB)

	hash := "hogehoge"
	expected := "https://en-courage.com"
	client.SetURL(hash, expected)
	got := client.GetURL(hash)
	if got != expected {
		t.Error("unexpected url, got: " + got + ", expected: " + expected)
	}

	client.DeleteURL(hash)
	deleted := client.GetURL(hash)
	if deleted != "" {
		t.Error("URL is not deleted, got: " + deleted)
	}
}
