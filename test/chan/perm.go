// errchk $G -e $D/$F.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var (
	cr <-chan int
	cs chan<- int
	c  chan int
)

func main() {
	cr = c  // ok
	cs = c  // ok
	c = cr  // ERROR "illegal types|incompatible|cannot"
	c = cs  // ERROR "illegal types|incompatible|cannot"
	cr = cs // ERROR "illegal types|incompatible|cannot"
	cs = cr // ERROR "illegal types|incompatible|cannot"

	c <- 0 // ok
	<-c    // ok
	//TODO(rsc): uncomment when this syntax is valid for receive+check closed
	//	x, ok := <-c	// ok
	//	_, _ = x, ok

	cr <- 0 // ERROR "send"
	<-cr    // ok
	//TODO(rsc): uncomment when this syntax is valid for receive+check closed
	//	x, ok = <-cr	// ok
	//	_, _ = x, ok

	cs <- 0 // ok
	<-cs    // ERROR "receive"
	////TODO(rsc): uncomment when this syntax is valid for receive+check closed
	////	x, ok = <-cs	// ERROR "receive"
	////	_, _ = x, ok

	select {
	case c <- 0: // ok
	case x := <-c: // ok
		_ = x

	case cr <- 0: // ERROR "send"
	case x := <-cr: // ok
		_ = x

	case cs <- 0: // ok
	case x := <-cs: // ERROR "receive"
		_ = x
	}
}
