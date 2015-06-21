package env

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/goanywhere/crypto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFindKeyValue(t *testing.T) {
	Convey("env.findKeyValue", t, func() {
		k, v := findKeyValue(" test= value")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")

		k, v = findKeyValue(" test= value")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")

		k, v = findKeyValue("\ttest=\tvalue\t\n")
		So(k, ShouldEqual, "test")
		So(v, ShouldEqual, "value")

		k, v = findKeyValue("export Test=\"Example\"")
		So(k, ShouldEqual, "Test")
		So(v, ShouldEqual, "Example")

		k, v = findKeyValue(`export Secret_Keys=IK-vyX7OuiftwyasT6NXnEYyPMj8fEDssJZdppKOs8Y4hZTtWfUILer73RbsG78Q`)
		So(k, ShouldEqual, "Secret_Keys")
		So(v, ShouldEqual, "IK-vyX7OuiftwyasT6NXnEYyPMj8fEDssJZdppKOs8Y4hZTtWfUILer73RbsG78Q")
	})
}

func TestLoad(t *testing.T) {
	filename := "/tmp/.env"
	// plain value without quote
	if dotenv, err := os.Create(filename); err == nil {
		defer dotenv.Close()
		defer os.Remove(filename)
		secret := crypto.Random(64)
		buffer := bufio.NewWriter(dotenv)
		buffer.WriteString(fmt.Sprintf("secret=%s\n", secret))
		buffer.WriteString("app=myapp\n")
		buffer.WriteString("export exportation=myexports")
		buffer.Flush()

		Convey("env.Load (without quoting)", t, func() {
			Set("root", "/tmp")
			Load("/tmp/.env")
			So(String("secret"), ShouldEqual, secret)
			So(String("app"), ShouldEqual, "myapp")
			So(String("exportation"), ShouldEqual, "myexports")
		})
	}
	// value with `` quote
	if dotenv, err := os.Create(filename); err == nil {
		defer dotenv.Close()
		defer os.Remove(filename)
		secret := crypto.Random(64)
		buffer := bufio.NewWriter(dotenv)
		buffer.WriteString(fmt.Sprintf("secret='%s'\n", secret))
		buffer.WriteString("app='myapp'\n")
		buffer.WriteString("export account='username'\n")
		buffer.Flush()

		Convey("env.Load (with quoting)", t, func() {
			Set("root", "/tmp")
			Load("/tmp/.env")
			So(String("secret"), ShouldEqual, secret)
			So(String("app"), ShouldEqual, "myapp")
			So(String("account"), ShouldEqual, "username")
		})
	}
	// value with `"` quote
	if dotenv, err := os.Create(filename); err == nil {
		defer dotenv.Close()
		defer os.Remove(filename)
		secret := crypto.Random(64)
		buffer := bufio.NewWriter(dotenv)
		buffer.WriteString(fmt.Sprintf("secret=\"%s\"\n", secret))
		buffer.WriteString("app=\"myapp\"\n")
		buffer.WriteString("export account=\"username\"\n")
		buffer.Flush()

		Convey("env.Load (with single quoting)", t, func() {
			Set("root", "/tmp")
			Load("/tmp/.env")
			So(String("secret"), ShouldEqual, secret)
			So(String("app"), ShouldEqual, "myapp")
			So(String("account"), ShouldEqual, "username")
		})
	}
}

func TestMap(t *testing.T) {
	type Person struct {
		Username  string
		Age       uint
		Kids      int
		Checked   bool
		Money     float64
		FirstName string `env:"FIRST_NAME"`
		Names     []string
	}
	var person Person
	Convey("env.Map", t, func() {
		Set("Username", "abc")
		Set("Age", 100)
		Set("Kids", 2)
		Set("Checked", true)
		Set("Money", 1234567890.0987654321)
		Set("FIRST_NAME", "abc")
		Set("Names", "a,b,c,d,e,f,g")
		Map(&person)

		So(person.Username, ShouldEqual, "abc")
		So(person.Age, ShouldEqual, 100)
		So(person.Kids, ShouldEqual, 2)
		So(person.Checked, ShouldBeTrue)
		So(person.Money, ShouldEqual, 1234567890.0987654321)
		So(person.FirstName, ShouldEqual, "abc")
		So(person.Names, ShouldResemble, []string{"a", "b", "c", "d", "e", "f", "g"})
	})
}

func TestString(t *testing.T) {
	Convey("env.String", t, func() {
		So(String("NotFound"), ShouldEqual, "")
		Set("Found", "something")
		So(String("Found"), ShouldEqual, "something")
		So(String("NotFound", "default"), ShouldEqual, "default")
	})
}

func TestStrings(t *testing.T) {
	Convey("env.Strings", t, func() {
		So(Strings("StringList"), ShouldBeNil)
		Set("StringList", "a,b,c")
		So(Strings("StringList"), ShouldResemble, []string{"a", "b", "c"})
		So(Strings("NotFound", []string{"a", "b", "c"}), ShouldResemble, []string{"a", "b", "c"})
	})
}

func TestInt(t *testing.T) {
	Convey("env.Int", t, func() {
		So(Int("integer"), ShouldEqual, 0)
		Set("integer", 123)
		So(Int("integer"), ShouldEqual, 123)
		So(Int("NotFound", 123), ShouldEqual, 123)
	})
}

func TestInt64(t *testing.T) {
	Convey("env.Int", t, func() {
		So(Int64("int64"), ShouldEqual, 0)
		Set("int64", 123)
		So(Int64("int64"), ShouldEqual, 123)
		So(Int64("NotFound", 123), ShouldEqual, 123)
	})
}

func TestUint(t *testing.T) {
	Convey("env.Uint", t, func() {
		So(Uint("uint"), ShouldEqual, 0)
		Set("uint", 123)
		So(Uint("uint"), ShouldEqual, 123)
		So(Uint("NotFound", 123), ShouldEqual, 123)
	})
}

func TestUint64(t *testing.T) {
	Convey("env.Uint64", t, func() {
		So(Uint64("uint64"), ShouldEqual, 0)
		Set("uint64", 123)
		So(Uint64("uint64"), ShouldEqual, 123)
		So(Uint64("NotFound", 123), ShouldEqual, 123)
	})
}

func TestBool(t *testing.T) {
	Convey("env.Bool", t, func() {
		So(Bool("bool"), ShouldBeFalse)
		Set("bool", true)
		So(Bool("bool"), ShouldBeTrue)
		So(Bool("NotFound", true), ShouldBeTrue)
	})
}

func TestFloat(t *testing.T) {
	Convey("env.Float", t, func() {
		So(Float("float64"), ShouldEqual, 0.0)
		Set("float64", 12345678990.0987654321)
		So(Float("float64"), ShouldEqual, 12345678990.0987654321)
		So(Float("NotFound", 12345678990.0987654321), ShouldEqual, 12345678990.0987654321)
	})
}
