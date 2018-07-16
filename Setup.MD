# Infrastructure Setup

All we got to on is essentially the following:

"Lastly he warned you that he has big​ plans for new features requiring both large additions to the game state as well as many new changes/additions to the API so you can make sure to have the system ready for these."

Meaning that while we have a very small glimpse into what the current requirements are, it may well take a complete turn around tomorrow and be nothing like it.

So that being said, we'll make a simple setup and try and get trough as much as possible with it, without getting too bogged down by the details.

You can find a small diagram of the main domain [here](Setup.jpg). It's the same you would get by googling "CQRS", as there really isn't much to draw a diagram over with the current level of information.

## Hosting

We're going to be using AWS EC2 instances, and all services are containerized in Docker.

## Main Domain

For the main domain we'll go with a CQRS setup, so there's space for the setup to grow into something worthy of the 8th wonder. For this very simple setup we'll just use a single type of the components: `query facade`, `thin data layer`, `command facade`,`command handler` and `event handler`. Adding in two busses, `command bus` and `event bus` along with two databases. A primary storage database and a materialized view optimized for the `thin data layer`.

Depending on the data model, we can use SQL databases (PostgreSQL, MySQL) or NoSQL (Cassandra, DynamoDB) for our storages.

## Authentication / Authorization

We're going to assume that there will be authentication in the system, if not used for the users we can use it for ourselves for internal accounts.

An example here could be an OAuth2 based Authentication, using encrypted JWT as bearer tokens which then contains a finely grained scope of access. This allows for decentralized authorization.

Authentication will use its own database, an ACID compliant one such as PostgreSQL.

## Client Access

With authentication and the main domain we got our client endpoint exposed and it can now access them. To allow scaling of the facade services, we'll use a load balancer to distribute the connections among them. As we're on amazon we can use Elastic Load Balancing (ELB), or if we want to do it ourselves Nginx is an option.

## Communication

You're most likely going to want to communicate with the user, outside of them sending requests to our services. To this end we'll set up some services, acting as gateways, to third party services handling things such as: E-mail, SMS and Push notifications.

These services will be available to communicate with all non-facade services in the setup.

## Internal Monitoring

Now we need to look at our own services and be sure that they are actually running, and running well enough.
Metrics (Services / Infrastructure) and Service health being gathered in a time series database (InfluxDB), error logs being aggregated in a text search database (ElasticSearch), and monitored by respective tools and send alarm to devs. We can here use our communication tools to send out alarms.

## Scaling

When we get alarmed that our system is being overloaded, it might already be too late to have a Dev spin up more resources, without it impeding performance.
Though I'd say this extends beyond the scope of 'simple', we can use technologies such as Kubernetes, Mesos and Terraform (+ Nomad) to automatically create more resources in the cloud.

## Internal tools

With the system being monitored and running, internal staff other than the Devs are interested in the system.

Admins and Moderators who needs to be able to take actions towards users of the system. Marketing or Financial staff who wants statistical insight into the system.

All of this requires its own setup, while it might not have the same requirement of scaling as the main domain, it can easily be the same size of complexity.