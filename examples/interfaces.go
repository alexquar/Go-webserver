/*
Go Interfaces - Notes and Examples

What is an Interface?
---------------------
An interface in Go is a type that specifies a set of method signatures (behavior), but does not provide implementation.
Any type that implements those methods satisfies the interface, implicitly.

Syntax:
--------
type InterfaceName interface {
	Method1(param1 type1, ...) returnType1
	Method2(param2 type2, ...) returnType2
}

Example: Basic Interface
------------------------
*/

// Define an interface
type Speaker interface {
	Speak() string
}

// Implement the interface with a struct
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + " says Woof!"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return c.Name + " says Meow!"
}

// Function that accepts the interface
func Announce(s Speaker) {
	println(s.Speak())
}

func main() {
	d := Dog{Name: "Buddy"}
	c := Cat{Name: "Whiskers"}

	Announce(d)
	Announce(c)
}

/*
Output:
Buddy says Woof!
Whiskers says Meow!

------------------------
Empty Interface
------------------------
The empty interface `interface{}` can hold values of any type (since all types implement zero methods).

Example:
*/

func PrintAnything(val interface{}) {
	println(val)
}

/*
------------------------
Type Assertion
------------------------
You can extract the concrete value from an interface using type assertion.

Example:
*/

func Describe(i interface{}) {
	s, ok := i.(string)
	if ok {
		println("It's a string:", s)
	} else {
		println("Not a string")
	}
}

/*
------------------------
Type Switch
------------------------
A type switch allows you to perform different actions based on the concrete type.

Example:
*/

func TypeSwitch(i interface{}) {
	switch v := i.(type) {
	case int:
		println("int:", v)
	case string:
		println("string:", v)
	default:
		println("unknown type")
	}
}

/*
------------------------
Interface Embedding
------------------------
Interfaces can embed other interfaces.

Example:
*/

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

/*
------------------------
Nil Interfaces
------------------------
A nil interface value has both its type and value set to nil.

Example:
*/

var i interface{}
println(i == nil) // true

/*
------------------------
Summary
------------------------
- Interfaces define behavior, not data.
- Satisfaction is implicit (no "implements" keyword).
- Use interfaces for abstraction and decoupling.
- The empty interface can hold any value.
- Use type assertions and type switches to extract concrete values.
*/