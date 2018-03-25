// Package natives provides native packages via a virtual filesystem.
//
// See documentation of parseAndAugment in github.com/gijit/gi/pkg/compiler/build.go
// for explanation of behavior used to augment the native packages using the files
// in src subfolder.
package natives

//go:generate vfsgendev -source="github.com/gijit/gi/pkg/compiler/natives".FS -tag=gopherjsdev
