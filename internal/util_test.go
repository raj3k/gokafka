package internal

import (
	"testing"
)

func BenchmarkIntToByteArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntToByteArray(12345) // Change the input value as needed
	}
}

func BenchmarkByteArrayToInt(b *testing.B) {
	data := IntToByteArray(12345) // Change the input value as needed

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ByteArrayToInt(data)
	}
}

func BenchmarkItoBSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ItoBSlice(12345) // Change the input value as needed
	}
}

func BenchmarkIntConversion(b *testing.B) {
	b.Run("IntToByteArray", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IntToByteArray(1234567890)
		}
	})

	b.Run("ItoBSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ItoBSlice(1234567890)
		}
	})
}
