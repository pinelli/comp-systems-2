package main

import (
	"fmt"
	"sync"
	"time"
)

type Instr struct {
	time time.Duration //ms
}

type Matrix struct {
	rows int
	cols int
}

func (p *Instr) execute() {
	<-time.NewTimer(p.time).C
}

func m1Planner(products chan *Instr, sums chan *Instr) {
	var group sync.WaitGroup
	group.Add(3)
	go m1Processor(products, sums, &group)
	go m1Processor(products, sums, &group)
	go m1Processor(products, sums, &group)
	group.Wait()
}

func m2Planner(products chan *Instr, sums chan *Instr) {
	m2Processor(products, sums)
}

func seqPlanner(products chan *Instr, sums chan *Instr) {
	var group sync.WaitGroup
	group.Add(1)
	go m1Processor(products, sums, &group)
	group.Wait()
}

func m1Processor(products chan *Instr, sums chan *Instr, group *sync.WaitGroup) {
	for {
		select {
		case instr := <-sums:
			instr.execute()
		case instr := <-products:
			instr.execute()
		default:
			(*group).Done()
			return
		}
	}
}

func m2Processor(products chan *Instr, sums chan *Instr) {
	cnt := 3
prods:
	for {
		select {
		case instr := <-products:
			cnt--
			if cnt == 0 {
				instr.execute()
				cnt = 3
			}
		default:
			break prods
		}
	}

	cnt = 3
	for {
		select {
		case instr := <-sums:
			cnt--
			if cnt == 0 {
				instr.execute()
				cnt = 3
			}
		default:
			return
		}
	}
}

func mkInstructions(m1 []Matrix, m2 []Matrix) (chan *Instr, chan *Instr) {
	prods := make(chan *Instr, 100)
	sums := make(chan *Instr, 100)
	for i := range m1 {
		m := m1[i].rows
		n := m1[i].cols
		k := m2[i].cols

		done := make(chan bool)
		for i := 0; i < m*k*n; i++ {
			go func() {
				prods <- &Instr{20 * time.Millisecond}
				done <- true
			}()
			<-done
		}

		for i := 0; i < m*k*(n-1); i++ {
			go func() {
				sums <- &Instr{10 * time.Millisecond}
				done <- true
			}()
			<-done
		}
	}
	return prods, sums
}

func main() {
	m1 := make([]Matrix, 10)
	m2 := make([]Matrix, 10)

	m1 = append(m1, Matrix{rows: 2, cols: 3})
	m2 = append(m1, Matrix{rows: 3, cols: 4})

	m1 = append(m1, Matrix{rows: 2, cols: 3})
	m2 = append(m1, Matrix{rows: 3, cols: 5})

	m1 = append(m1, Matrix{rows: 2, cols: 5})
	m2 = append(m1, Matrix{rows: 5, cols: 3})

	prods, sums := mkInstructions(m1, m2)

	start := time.Now()
	fmt.Println("Execute...")

	//Uncomment needed algorithm
	m1Planner(prods, sums)
	//m2Planner(prods, sums)
	//seqPlanner(prods, sums)

	fmt.Println("Execution time:", time.Since(start))
}
