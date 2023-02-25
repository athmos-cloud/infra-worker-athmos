package workdir

type Workdir struct {
	Files   []File   `bson:"files"`
	Folders []Folder `bson:"folders"`
}

func (w *Workdir) AddFile(file File) {
	w.Files = append(w.Files, file)
}
func (w *Workdir) AddFolder(folder Folder) {
	w.Folders = append(w.Folders, folder)
}
