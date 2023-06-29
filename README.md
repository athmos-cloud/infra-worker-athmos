# Athmos Infra Worker

This application is dedicated to provision infrastructure on the cloud. AWS and GCP are
currently supported.

This application handle requests  both synchronously throught an HTTP server and asynchronously
by listening on a RabbitMQ.

You will find in the `config.yaml` file a defautl configuration for the application.

Moreover, you need to set at least the following environment variables :

* `CONFIG_FILE_LOCATION` which is the path to the configuration file.
* `REDIRECTION_URL` which is the URL to redirect to get resource individually.
* `KUBECONFIG_PATH` which is the path to the kubeconfig file.
* `PORT` which is the port to listen to.
* `MONGO_ADDRESS` which is the address of the MongoDB instance.
* `MONGO_DATABASE` which is the name of the MongoDB database.
* `MONGO_USERNAME` which is the username to connect to the MongoDB instance.
* `MONGO_PASSWORD` which is the password to connect to the MongoDB instance.
* `RABBITMQ_ADDRESS` which is the address of the RabbitMQ instance.
* `RABBITMQ_USERNAME` which is the username to connect to the RabbitMQ instance.
* `RABBITMQ_PASSWORD` which is the password to connect to the RabbitMQ instance.
* `RABBITMQ_PORT` which is the port to connect to the RabbitMQ instance.


This application relies on the [crossplane](https://crossplane.io/) operator. Hence, you will
have to get a valid kubeconfig file (at the specified path) to connect to the cluster where the operator is running.
You will also need to install the operators for GCP and AWS providers. You will find the configs
into the `config` folder.
