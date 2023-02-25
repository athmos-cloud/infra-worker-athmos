package workdir

type File struct {
	Name    string `bson:"name"`
	Content string `bson:"content"`
}
