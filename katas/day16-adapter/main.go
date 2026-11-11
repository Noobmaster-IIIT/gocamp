// Command adapter shows how a plain function can satisfy a one-method interface
// via a named function type with a method (the http.HandlerFunc pattern).
//
// Run it from the gocamp root:  go run ./katas/day16-adapter
package main

import "fmt"

// Greeter is an ordinary one-method interface.
type Greeter interface {
	Greet(name string) string
}

// ---------------------------------------------------------------------------
// WAY 1: the usual way, a struct type with a Greet method.
// ---------------------------------------------------------------------------

// PoliteGreeter satisfies Greeter because it has Greet(string) string.
type PoliteGreeter struct{}

func (p PoliteGreeter) Greet(name string) string {
	return "Good evening, " + name
}

// ---------------------------------------------------------------------------
// WAY 2: the adapter, a named FUNCTION type with a Greet method.
// Written ONCE; after that, any func(string) string can become a Greeter.
// ---------------------------------------------------------------------------

// GreeterFunc is a named type whose underlying type is func(string) string.
// A value of this type IS a function.
type GreeterFunc func(name string) string

// Greet belongs to the TYPE GreeterFunc. The receiver f is just the first
// argument, and here it happens to be a function, so "using" it means CALLING it.
func (f GreeterFunc) Greet(name string) string {
	fmt.Println("   ④ inside GreeterFunc.Greet: about to call f(name)")
	return f(name) // ⑤ f is the function stored in the interface, e.g. polite
}

// welcome only knows about the Greeter interface. It has no idea whether it got
// a struct or a function.
func welcome(g Greeter) {
	fmt.Printf("   ③ welcome got g with dynamic type %T\n", g)
	fmt.Println("   ⑥", g.Greet("Abhirup"))
}

func main() {
	fmt.Println("WAY 1: struct")
	welcome(PoliteGreeter{})

	fmt.Println("\nWAY 2: function + adapter")
	// ① a plain function value. Its type, func(string) string, has NO methods.
	polite := func(n string) string { return "Good evening, " + n }

	// welcome(polite)  // ✗ compile error: func(string) string does not implement
	//                  //   Greeter (missing method Greet). Uncomment to see it.

	// ② CONVERSION (like Celsius(21.5)), not a call: same function, new type label.
	//    The type GreeterFunc has a Greet method, so the value now satisfies Greeter.
	g := GreeterFunc(polite)
	welcome(g)

	fmt.Println("\nWAY 2 again, inline (step ① happens right inside the conversion)")
	welcome(GreeterFunc(func(n string) string { return "yo " + n }))

	fmt.Println("\nBonus: the receiver really is just the first argument")
	fmt.Println("   ", GreeterFunc.Greet(g, "Abhirup")) // method expression: Type.Method(receiver, args...)
}
