package tick

// Template definition for TICK scripts
type Template struct {
	ID       string
	Template string
}

// Templates contains the global list of available templates
var Templates = map[string]Template{
	"autoscaling": AutoScaling,
}

// AutoScaling TICK template can be used to calculate against cpu/memory/queue sizes as long as the data is
// aggregated to InfluxDB.  The easiest way to funnel data to Influx is to use Telegraf
var AutoScaling = Template{
	ID: "autoscaling",
	Template: `
// Database
var database string
// Retention policy for database
var retentionPolicy string
// Dataset collected within the retention policy
var measurement string
// Optional where filter
var whereFilter = lambda: TRUE
// Field data to use for the processing
var field string
// The time scale to calculate the average against
var movingAverageCount = 60
// Namespace the deployment lives in
var namespace = 'default'
// Deployment this is scaling for
var deploymentName = 'placeholder'
// Threshold for triggering
var target = 10.0
// Time interval per scaling up
var scalingCooldown = 2m
// Time interval per scaling down
var descalingCooldown = 5m
// Minimum replica count to maintain regardless of needs
var minReplicaCount = 1
// Maximum replica count to stop scaling at regardless of needs
var maxReplicaCount = 20
	
stream
	|from()
		.database(database)
		.retentionPolicy(retentionPolicy)
		.measurement(measurement)
		.where(whereFilter)
		.truncate(1s)
	|movingAverage(field, movingAverageCount)
		.as('averageResource')	
	|eval(lambda: int(ceil("averageResource" / float(target))))
		.as('replicaCount')
	|k8sAutoscale()
		// We are scaling a deployment.
		.kind('deployments')
		// The namespace of the deployment
		.namespace(namespace)
		// The name of the replicaset to scale is found in the 'replicaset' tag.
		.resourceName(deploymentName)
		// Set the cool down timer values.
		.increaseCooldown(scalingCooldown)
		.decreaseCooldown(descalingCooldown)
		// The minimum amount of replica to have regardless of averages
		.min(minReplicaCount)
		// The maximum amount of replica to have regardless of averages
		.max(maxReplicaCount)
		// Compute the desired number of replicas based on the
		// avg_point_per_second and target values.
		.replicas(lambda: "replicaCount")
	`,
}
