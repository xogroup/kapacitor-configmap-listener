package tick

// Template definition for TICK scripts
type Template struct {
	ID       string
	Template string
}

// Global list of templates
var Templates = map[string]Template{
	"autoscaling": AutoScaling,
}

// AutoScaling TICK tempate
var AutoScaling = Template{
	ID: "autoscaling",
	Template: `
// database
var database string
// retention policy for database
var retentionPolicy string
// dataset collected within the retention policy
var measurement string
// Optional where filter
var where_filter = lambda: TRUE
// Optional list of group by dimensions
var groups = ['host']
// Which field to process
var field string
// The time scale to calculate the average against
var movingAverageCount = 60
// Deployment this is scaling for
var deploymentName = 'placeholder'
// Threshold for triggering
var target = 10.0
// time interval per scaling up
var scalingCooldown = 1m
// time interval per scaling down
var descalingCooldown = 2m
	
stream
	|from()
		.database(database)
		.retentionPolicy(retentionPolicy)
		.measurement(measurement)
		.where(where_filter)
		.groupBy(groups)
		.truncate(1s)
	// Compute the rate of requests per second per host
	|derivative(field)
		.as('point_per_second')
		.unit(1s)
		.nonNegative()
	|alert()
		.crit(lambda: "point_per_second" > target)
		.log('/var/log/test.log')
	|sum('point_per_second')
		.as('total_point_per_second')
	|movingAverage('total_point_per_second', movingAverageCount)
		.as('avg_point_per_second')
	// add window()
	|k8sAutoscale()
		// We are scaling a deployment.
		.kind('deployments')
		// The name of the replicaset to scale is found in the 'replicaset' tag.
		.resourceName(deploymentName)
		// Set the cool down timer values.
		.increaseCooldown(scalingCooldown)
		.decreaseCooldown(descalingCooldown)
		// Compute the desired number of replicas based on the
		// avg_point_per_second and target values.
		.replicas(lambda: int(ceil("avg_point_per_second" / target)))
	`,
}

// stream
// |from()
// 	.database(‘telegraf’)
// 	.retentionPolicy(’2-weeks’)
// 	.measurement('docker_container_cpu')
// 	.where(lambda: "io.kubernetes.container.name" == 'knot-proxy-qa')
// 	.groupBy('host')
// |movingAverage('usage_percent', moving_avg_count)
// 	.as('data_per_second')
// |alert()
// 	.crit(lambda: "data_per_second" > target)
// 	.log('/var/log/test.log')
// |k8sAutoscale()
// 	// We are scaling a replicaset.
// 	.kind('deployments')
// 	// The name of the replicaset to scale is found in the 'replicaset' tag.
// 	.resourceName('nginx-app')
// 	// Set the cool down timer values.
// 	.increaseCooldown(increase_cooldown)
// 	.decreaseCooldown(decrease_cooldown)
// 	// Compute the desired number of replicas based on the
// 	// avg_requests_per_second and target values.
// 	.replicas(lambda: int(ceil("data_per_second" / target)))
