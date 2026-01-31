# Match Runner

The Match Runner is a process, that will pull the next Game from a RabbitMQ instance, when he is ready.
He will pull the required docker images and run teh match on the single shot server instance.
Afterwards he's going to publish the result back into a RabbitMQ.
