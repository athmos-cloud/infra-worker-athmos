package workdir

type Folder struct {
	Name    string   `bson:"name"`
	Files   []File   `bson:"files"`
	Folders []Folder `bson:"folders"`
}

func (f *Folder) AddFile(file File) {
	f.Files = append(f.Files, file)
}
func (f *Folder) AddFolder(folder Folder) {
	f.Folders = append(f.Folders, folder)
}
