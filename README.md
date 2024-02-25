# loadbalancer

This is a simple implementation of a loadbalancer using Go and its std libraries

I have used the round robin method to distribute the load from the load balancer to the multiple servers available in the background

I have used the std library provided reverse proxy functionality tp forward the request from the load balancer to the actual servers.

In the current implementation, it only supports roundrobin. I plan to enhance it with few other methods like random, leastconnectio etc in the upcoming weeks.

While running the module, you may use the command line argument alg to pass the load balancing method. It will default to robin if not set.

