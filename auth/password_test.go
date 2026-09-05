package auth

import "testing"

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	password := "secret"

	hash1 := HashPassword(password)
	hash2 := HashPassword(password)

	if hash1 == hash2 {
		t.Fatal("hashes are identical, expected different salts")
	}
}

func TestVerifySamePassword(t *testing.T) {
	password1 := "secret"
	password2 := "secret"
	hashed := HashPassword(password1)

	ok := VerifyPassword(password2, hashed)
	if !ok {
		t.Fatal("same password verification failed, expected success")
	}
}

func TestVerifyDifferentPassword(t *testing.T) {
	password1 := "secret"
	password2 := "wrong"
	hashed := HashPassword(password1)

	ok := VerifyPassword(password2, hashed)
	if ok {
		t.Fatal("different password verification succeeded, expected failure")
	}
}

func TestHashPasswordPHCFormat(t *testing.T) {
	password := "secret"
	hashed := HashPassword(password)

	salt, hash, params := decodeHash(hashed)
	if salt == nil || hash == nil || params == nil {
		t.Fatal("hash decoding failed")
	}

	if len(salt) != int(saltLength) {
		t.Fatalf("unexpected salt length: got %d, want %d", len(salt), saltLength)
	}

	if len(hash) != int(keyLength) {
		t.Fatalf("unexpected hash length: got %d, want %d", len(hash), keyLength)
	}

	if params.memory != memory {
		t.Fatalf("unexpected memory: got %d, want %d", params.memory, memory)
	}

	if params.iterations != iterations {
		t.Fatalf("unexpected iterations: got %d, want %d", params.iterations, iterations)
	}

	if params.parallelism != parallelism {
		t.Fatalf("unexpected parallelism: got %d, want %d", params.parallelism, parallelism)
	}
}
