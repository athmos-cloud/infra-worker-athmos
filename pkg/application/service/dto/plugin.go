package dto

type CreatePluginInstanceRequest struct {
	PluginId  string            `bson:"id"`
	ProjectId string            `bson:"projectId"`
	Values    map[string]string `bson:"values"`
}

type UpdatePluginRequest struct {
	PluginInstanceId string            `bson:"id"`
	ProjectId        string            `bson:"projectId"`
	Values           map[string]string `bson:"values"`
}

type GetPluginInstanceRequest struct {
	PluginInstanceId string `bson:"id"`
	ProjectId        string `bson:"projectId"`
}

type DeletePluginInstanceRequest struct {
	PluginInstanceId string `bson:"id"`
	ProjectId        string `bson:"projectId"`
}
