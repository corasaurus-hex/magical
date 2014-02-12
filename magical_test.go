package main

import (
	"math/rand"
	"testing"
)

func BenchmarkGenerateIds10(b *testing.B) {
	setup()
	n := 10
	quit := make(chan bool, n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for x := 0; x < n; x++ {
			go func() {
				for d := 0; d < n; d++ {
					generateIds(rand.Intn(300) + 1)
				}
				quit <- true
			}()
		}

		for x := 0; x < n; x++ {
			<- quit
		}
	}
}

func BenchmarkGenerateIds20(b *testing.B) {
	setup()
	n := 20
	quit := make(chan bool, n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for x := 0; x < n; x++ {
			go func() {
				for d := 0; d < n; d++ {
					generateIds(rand.Intn(300) + 1)
				}
				quit <- true
			}()
		}

		for x := 0; x < n; x++ {
			<- quit
		}
	}
}

func BenchmarkGenerateIds30(b *testing.B) {
	setup()
	n := 30
	quit := make(chan bool, n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for x := 0; x < n; x++ {
			go func() {
				for d := 0; d < n; d++ {
					generateIds(rand.Intn(300) + 1)
				}
				quit <- true
			}()
		}

		for x := 0; x < n; x++ {
			<- quit
		}
	}
}

func BenchmarkGenerateIds40(b *testing.B) {

	setup()
	n := 40
	quit := make(chan bool, n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for x := 0; x < n; x++ {
			go func() {
				for d := 0; d < n; d++ {
					generateIds(rand.Intn(300) + 1)
				}
				quit <- true
			}()
		}

		for x := 0; x < n; x++ {
			<- quit
		}
	}
}
