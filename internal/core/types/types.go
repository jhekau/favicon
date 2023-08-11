package types

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 * TDD
 */


//
// type FileName string
// func (f FileName) String() string {
// 	return string(f)
// }

type FilePath string
func (f FilePath) String() string {
	return string(f)
}

// type Folder string
// func (f Folder) String() string {
// 	return string(f)
// }

//
// type FileExists int
// const (
// 	FileExistsNotCheck FileExists = 0
// 	FileExistsOK FileExists = 1
// 	FileExistsNOT FileExists = -1
// )



//
type URLHref string
func (n URLHref) String() string {
	return string(n)
}