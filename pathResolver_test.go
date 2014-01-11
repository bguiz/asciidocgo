package asciidocgo

import (
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPathResolver(t *testing.T) {

	Convey("A pathResolver can be initialized", t, func() {

		Convey("By default, a pathResolver can be created", func() {
			So(newPathResolver(0, ""), ShouldNotBeNil)
		})
		Convey("By default, a pathResolver has a system path separator", func() {
			So(newPathResolver(0, "").FileSeparator(), ShouldEqual, os.PathSeparator)
			So(newPathResolver('/', "").FileSeparator(), ShouldNotEqual, os.PathSeparator)
			So(newPathResolver('/', "").FileSeparator(), ShouldEqual, '/')
		})

		Convey("By default, a pathResolver has a current working path", func() {
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			So(newPathResolver(0, "").WorkingDir(), ShouldEqual, pwd)
			So(newPathResolver(0, "C:\\").WorkingDir(), ShouldEqual, "C:\\")
			So(newPathResolver(0, "test").WorkingDir(), ShouldEqual, pwd+string(os.PathSeparator)+"test")
			//So(newPathResolver(0, "panicnoroot").WorkingDir(), ShouldEqual, pwd)

		})
		Convey("A pathResolver should not panic on getting pwd", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = newPathResolver(0, "panic on os.Getwd")
		})
		Convey("A pathResolver should not panic on filepath.Abs", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = newPathResolver(0, "panic on filepath.Abs")
		})
	})

	Convey("A pathResolver can test for a web path", t, func() {
		So(IsWebRoot(""), ShouldBeFalse)
		So(IsWebRoot("a"), ShouldBeFalse)
		So(IsWebRoot("\\a\\b/c"), ShouldBeFalse)
		So(IsWebRoot("/a/b/c"), ShouldBeTrue)
	})

	Convey("A pathResolver can replace backslash by slash", t, func() {
		So(Posixfy(""), ShouldEqual, "")
		So(Posixfy("a/b/c"), ShouldEqual, "a/b/c")
		So(Posixfy("a\\b\\c"), ShouldEqual, "a/b/c")
	})

	Convey("A pathResolver can test for root", t, func() {
		Convey("A Path starting with C:/ is root", func() {
			So(IsRoot(""), ShouldBeFalse)
			So(IsRoot("C:\\"), ShouldBeTrue)
			So(IsRoot("C:/"), ShouldBeTrue)
			So(IsRoot("C:\\a/b/"), ShouldBeTrue)
		})
	})

	Convey("A pathResolver can test for web root", t, func() {
		Convey("A Path starting with / is web root", func() {
			So(IsWebRoot(""), ShouldBeFalse)
			So(IsWebRoot("C:\\"), ShouldBeFalse)
			So(IsWebRoot("\\"), ShouldBeFalse)
			So(IsWebRoot("/"), ShouldBeTrue)
			So(IsWebRoot("/a/b/"), ShouldBeTrue)
		})
	})

	Convey("A pathResolver can partition a path", t, func() {
		pathSegments, root, posixPath := PartitionPath("", false)
		So(len(pathSegments), ShouldEqual, 0)
		So(root, ShouldEqual, "")
		So(posixPath, ShouldEqual, "")

		Convey("A Path starting with dot has a dot root", func() {
			pathSegments, root, posixPath := PartitionPath(".", false)
			So(len(pathSegments), ShouldEqual, 0)
			So(root, ShouldEqual, ".")
			So(posixPath, ShouldEqual, ".")

			pathSegments, root, posixPath = PartitionPath(".\\a/b", false)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, ".")
			So(posixPath, ShouldEqual, "./a/b")

		})
		Convey("A Partition removes self-reference path", func() {
			pathSegments, root, posixPath := PartitionPath("a\\b/./c", false)
			So(len(pathSegments), ShouldEqual, 3)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b/./c")
			pathSegments, root, posixPath = PartitionPath("C:/a\\b/./c", false)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, "C:")
			So(posixPath, ShouldEqual, "C:/a/b/./c")
			pathSegments, root, posixPath = PartitionPath("/a\\b/./c", true)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "/a/b/./c")
		})
	})
}
