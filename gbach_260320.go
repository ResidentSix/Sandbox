// goldbach_prop.go
// Verify Goldbach via symmetric propagation: try k=1, then k=2 from previously proven levels.
// Uses a simple sieve; runs in O(N log log N) time and O(N) memory.
// Example: go run goldbach_prop.go 1000000   // checks up to 2,000,000

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// ---- read limit -------------------------------------------------
	var limit int
	if len(os.Args) > 1 {
		fmt.Sscan(os.Args[1], &limit)
	} else {
		fmt.Print("Enter maximum n (check even numbers up to 2n): ")
		fmt.Fscanln(bufio.NewReader(os.Stdin), &limit)
	}
	if limit < 2 {
		fmt.Println("limit must be at least 2")
		return
	}

	// ---- sieve of Eratosthenes (byte array for speed) ---------------
	maxVal := 2*limit + 10 // we need to test numbers up to 2n
	isPrime := make([]bool, maxVal+1)
	for i := 2; i <= maxVal; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= maxVal; p++ {
		if isPrime[p] {
			for step := p * p; step <= maxVal; step += p {
				isPrime[step] = false
			}
		}
	}

	// ---- store Goldbach pairs for each level n -----------------------
	pairs := make([][][2]int, limit+1) // pairs[n] = list of (p,q) with p+q=2n, p≤q

	// ---- main loop --------------------------------------------------
	start := timeNow()
	k1Count, k2Count := 0, 0
	var firstFail int

	for n := 2; n <= limit; n++ {
		// collect pairs for this n
		var lst [][2]int
		for p := 2; p <= n; p++ {
			q := 2*n - p
			if isPrime[p] && isPrime[q] {
				lst = append(lst, [2]int{p, q})
			}
		}
		pairs[n] = lst

		if n < 2 { // n=0,1 not needed
			continue
		}
		if n < 4 { // 2n = 4,6 are trivial bases (2+2, 3+3)
			continue
		}

		// try k=1 first
		found := false
		if n-1 >= 2 {
			for _, pr := range pairs[n-1] {
				p, q := pr[0], pr[1]
				// (p+2, q)
				if isPrime[p+2] && isPrime[q] {
					found = true
					break
				}
				// (p, q+2)
				if isPrime[p] && isPrime[q+2] {
					found = true
					break
				}
			}
		}
		if found {
			k1Count++
			continue
		}

		// try k=2 (only needed when k=1 failed)
		if n-2 >= 2 {
			for _, pr := range pairs[n-2] {
				p, q := pr[0], pr[1]
				// (p+4, q)
				if isPrime[p+4] && isPrime[q] {
					found = true
					break
				}
				// (p+2, q+2)
				if isPrime[p+2] && isPrime[q+2] {
					found = true
					break
				}
				// (p, q+4)
				if isPrime[p] && isPrime[q+4] {
					found = true
					break
				}
			}
		}
		if found {
			k2Count++
			continue
		}

		// no k worked → would be a counter‑example
		firstFail = 2*n
		break
	}

	elapsed := timeSince(start)

	if firstFail > 0 {
		fmt.Printf("\n!! COUNTEREXAMPLE FOUND at 2n = %d !!\n", firstFail)
		os.Exit(1)
	}

	fmt.Printf("\nVerification complete for even numbers up to 2n = %d\n", 2*limit)
	fmt.Printf("  Total numbers checked : %d\n", limit-1) // from n=2 to limit
	fmt.Printf("  Resolved by k=1       : %d (%.1f%%)\n", k1Count, 100.0*float64(k1Count)/float64(limit-1))
	fmt.Printf("  Resolved by k=2       : %d (%.1f%%)\n", k2Count, 100.0*float64(k2Count)/float64(limit-1))
	fmt.Printf("  Elapsed time          : %.3f s\n", elapsed.Seconds())
}

// helpers for timing (using time package)
import "time"

func timeNow() time.Time { return time.Now() }
func timeSince(t time.Time) time.Duration { return time.Since(t) }
