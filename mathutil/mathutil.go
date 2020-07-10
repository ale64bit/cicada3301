package mathutil

var (
	primes []int
)

func Pow(n, e int) int {
	switch {
	case e <= 0:
		return 1
	case e%2 == 0:
		x := Pow(n, e>>1)
		return x * x
	default:
		return n * Pow(n, e-1)
	}
}

func ModPow(n, e, m int) int {
	switch {
	case e <= 0:
		return 1
	case e%2 == 0:
		x := ModPow(n, e>>1, m)
		return (x * x) % m
	default:
		return (n * ModPow(n, e-1, m)) % m
	}
}

func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func Coprime(a, b int) bool {
	return GCD(a, b) == 1
}

func Phi(n int) int {
	factors := Factor(n)
	unique := map[int]bool{}
	for _, p := range factors {
		unique[p] = true
	}
	ret := n
	for p := range unique {
		ret /= p
		ret *= (p - 1)
	}
	return ret
}

func Mu(n int) int {
	factors := Factor(n)
	unique := map[int]bool{}
	for _, p := range factors {
		if unique[p] {
			return 0
		}
		unique[p] = true
	}
	if len(unique)%2 == 0 {
		return 1
	}
	return -1
}

func Factor(n int) []int {
	var factors []int
	i := 0
	for n > 1 {
		if i >= len(primes) {
			panic("factor is too large")
		}
		for n%primes[i] == 0 {
			n /= primes[i]
			factors = append(factors, primes[i])
		}
		i++
	}
	return factors
}

func PrimeAt(k int) int {
	if k >= len(primes) {
		panic("too large prime")
	}
	return primes[k]
}

func Primes(n int) []int {
	if n < 3 {
		return []int{2}
	}
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0] = false
	isPrime[1] = false
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i + i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	var primes []int
	for i, ok := range isPrime {
		if ok {
			primes = append(primes, i)
		}
	}
	return primes
}

func init() {
	primes = Primes(1000000)
}
