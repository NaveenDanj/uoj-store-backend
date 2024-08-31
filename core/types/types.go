package types

import "io/fs"

type FileChunkMetaData struct {
	Sequence     []string
	FileInfo     fs.FileInfo
	Extension    string
	OriginalName string
}

type RSA struct {
	PubliKeyPath   string
	PrivateKeyPath string
}
