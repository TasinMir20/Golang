package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	var i int32 = 0
	app.Get("/", func(c fiber.Ctx) error {

		atomic.AddInt32(&i, 1)

		fmt.Println(i)

		return c.JSON(fiber.Map{
			"message": "Hello from Go",
			"request": i,
		})
	})

	// Define another route
	app.Get("/greet/:name", func(c fiber.Ctx) error {
		name := c.Params("name")
		return c.JSON(fiber.Map{
			"greeting": "Hello, " + name + "!",
		})
	})

	// CPU intensive route
	app.Get("/fibonacci/:n", func(c fiber.Ctx) error {
		n, err := strconv.Atoi(c.Params("n"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Please provide a valid number",
			})
		}

		startTime := time.Now()

		// Recursive Fibonacci calculation (intentionally CPU intensive)
		callCount := 0
		var fib func(int) int
		fib = func(n int) int {
			callCount++
			if callCount%1000000 == 0 { // Log every millionth call
				fmt.Printf("⚡ Made %d million recursive calls...\n", callCount/1000000)
			}

			if n <= 1 {
				return n
			}
			return fib(n-1) + fib(n-2)
		}

		result := fib(n)
		duration := time.Since(startTime)
		minutes := int(duration.Minutes())
		seconds := int(duration.Seconds()) % 60
		milliseconds := int(duration.Milliseconds()) % 1000
		timeTook := fmt.Sprintf("%d minutes %d seconds %d milliseconds", minutes, seconds, milliseconds)

		return c.JSON(fiber.Map{
			"number":   n,
			"result":   result,
			"timeTook": timeTook,
		})
	})

	app.Get("/stress/:seconds", func(c fiber.Ctx) error {
		seconds, err := strconv.Atoi(c.Params("seconds"))
		if err != nil || seconds <= 0 || seconds > 30 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Please provide a valid number of seconds (1-30)",
			})
		}

		numCPU := runtime.NumCPU()
		var wg sync.WaitGroup
		startTime := time.Now()

		// Atomic counters for logging
		var totalExponentialOps uint64
		var totalLinearOps uint64
		var totalMemoryOps uint64

		fmt.Printf("🚀 Starting stress test for %d seconds on %d CPU cores\n", seconds, numCPU)

		// Global memory pressure - keeps memory allocated throughout the test
		globalMemory := make([][]byte, 0, numCPU*2)

		for i := 0; i < numCPU; i++ {
			wg.Add(1)
			go func(coreID int) {
				defer wg.Done()
				localExpOps := uint64(0)
				localLinearOps := uint64(0)
				localMemoryOps := uint64(0)

				fmt.Printf("⚡ Core %d: Starting intensive calculations\n", coreID)

				// Memory chunks for this core
				const memoryChunkSize = 256 * 1024 * 1024 // 256MB per chunk
				localMemory := make([][]byte, 0, 4)

				for time.Since(startTime).Seconds() < float64(seconds) {
					// Memory pressure operations
					chunk := make([]byte, memoryChunkSize)
					for i := range chunk[:1024] { // Initialize first 1KB for actual memory allocation
						chunk[i] = byte(i * coreID)
					}
					localMemory = append(localMemory, chunk)
					if len(localMemory) > 2 { // Keep ~512MB per core
						localMemory = localMemory[1:]
					}
					localMemoryOps++

					// Original CPU-intensive exponential calculations
					for i := 0; i < 30; i++ {
						iterations := 1 << i
						if i%5 == 0 {
							fmt.Printf("🔄 Core %d: Processing 2^%d = %d iterations\n", coreID, i, iterations)
						}

						result := 0
						for j := 0; j < iterations; j++ {
							x := float64(j)
							result += int(math.Pow(x, x) +
								math.Sqrt(x) +
								math.Sin(x) +
								math.Cos(x) +
								math.Tan(x) +
								math.Log(x+1) +
								math.Exp(math.Mod(x, 10)))
						}
						localExpOps += uint64(iterations)
					}

					// Additional intensive calculations
					for i := 0; i < 1000000; i++ {
						x := float64(i)
						math.Pow(x, x)
						math.Sqrt(x)
						math.Sin(x)
						math.Cos(x)
						math.Tan(x)
						math.Log(x + 1)
						math.Exp(math.Mod(x, 10))
					}
					localLinearOps += 1000000

					// Periodically add to global memory pressure
					if localMemoryOps%10 == 0 {
						globalChunk := make([]byte, 128*1024*1024) // 128MB
						for i := range globalChunk[:1024] {
							globalChunk[i] = byte(i * coreID)
						}
						globalMemory = append(globalMemory, globalChunk)
						if len(globalMemory) > numCPU*2 {
							globalMemory = globalMemory[1:]
						}
					}

					// Log progress every second
					if int(time.Since(startTime).Seconds())%1 == 0 {
						atomic.AddUint64(&totalExponentialOps, localExpOps)
						atomic.AddUint64(&totalLinearOps, localLinearOps)
						atomic.AddUint64(&totalMemoryOps, localMemoryOps)
						fmt.Printf("📊 Core %d: Processed %d exp ops, %d linear ops, %d memory ops (Memory: %d chunks)\n",
							coreID, localExpOps, localLinearOps, localMemoryOps, len(localMemory))
						localExpOps = 0
						localLinearOps = 0
						localMemoryOps = 0
					}
				}

				// Cleanup local memory
				localMemory = nil
				fmt.Printf("✅ Core %d: Completed calculations\n", coreID)
			}(i)
		}

		wg.Wait()
		// Clear global memory
		globalMemory = nil
		duration := time.Since(startTime)

		finalStats := fiber.Map{
			"message":          fmt.Sprintf("Stressed %d CPU cores for %v", numCPU, duration.Round(time.Millisecond)),
			"cores":            numCPU,
			"total_exp_ops":    totalExponentialOps,
			"total_linear_ops": totalLinearOps,
			"total_memory_ops": totalMemoryOps,
			"duration_seconds": duration.Seconds(),
		}

		fmt.Printf("🏁 Stress test completed:\n")
		fmt.Printf("   - Duration: %v\n", duration.Round(time.Millisecond))
		fmt.Printf("   - Total exponential operations: %d\n", totalExponentialOps)
		fmt.Printf("   - Total linear operations: %d\n", totalLinearOps)
		fmt.Printf("   - Total memory operations: %d\n", totalMemoryOps)
		fmt.Printf("   - Average ops per second: %d\n", (totalExponentialOps+totalLinearOps+totalMemoryOps)/uint64(seconds))

		return c.JSON(finalStats)
	})

	// Start the server on port 3000
	if err := app.Listen(":4000"); err != nil {
		panic(err)
	}

}
