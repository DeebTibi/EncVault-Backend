package Server

import "fmt"

var serviceToIps = make(map[string][]string)

// RegisterService registers a service with the registry
func RegisterService(serviceName, ip string) {
	fmt.Printf("Registering service %v with IP %v\n", serviceName, ip)
	if _, ok := serviceToIps[serviceName]; !ok {
		serviceToIps[serviceName] = []string{}
	}
	serviceToIps[serviceName] = append(serviceToIps[serviceName], ip)
}

// GetServiceIps returns the list of IPs for a service
func GetServiceIps(serviceName string) []string {
	ips, ok := serviceToIps[serviceName]
	fmt.Printf("Getting IPs for service %v: %v\n", serviceName, ips)
	if !ok {
		return nil
	}
	return ips
}

// UnregisterService unregisters a service from the registry
func UnregisterService(serviceName, ip string) {
	fmt.Printf("Unregistering service %v with IP %v\n", serviceName, ip)
	if _, ok := serviceToIps[serviceName]; !ok {
		return
	}

	ips := serviceToIps[serviceName]
	for i, v := range ips {
		if v == ip {
			serviceToIps[serviceName] = append(ips[:i], ips[i+1:]...)
			break
		}
	}
	ips = serviceToIps[serviceName]
	if len(ips) == 0 {
		delete(serviceToIps, serviceName)
	}
}
