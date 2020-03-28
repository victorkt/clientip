package clientip

import (
	"net"
	"net/http"
	"strings"
)

// FromRequest returns the client IP address from the HTTP request
func FromRequest(r *http.Request) net.IP {
	// Standard headers used by Amazon EC2, Heroku, and others.
	if ip := net.ParseIP(r.Header.Get("x-client-ip")); ip != nil {
		return ip
	}

	// Load-balancers (AWS ELB) or proxies.
	if ip := fromXForwardedFor(r.Header.Get("x-forwarded-for")); ip != nil {
		return ip
	}

	// Cloudflare.
	// @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
	// CF-Connecting-IP - applied to every request to the origin.
	if ip := net.ParseIP(r.Header.Get("cf-connecting-ip")); ip != nil {
		return ip
	}

	// Fastly and Firebase hosting header (When forwared to cloud function)
	if ip := net.ParseIP(r.Header.Get("fastly-client-ip")); ip != nil {
		return ip
	}

	// Akamai and Cloudflare: True-Client-IP.
	if ip := net.ParseIP(r.Header.Get("true-client-ip")); ip != nil {
		return ip
	}

	// Default nginx proxy/fcgi; alternative to x-forwarded-for, used by some proxies.
	if ip := net.ParseIP(r.Header.Get("x-real-ip")); ip != nil {
		return ip
	}

	// (Rackspace LB and Riverbed's Stingray)
	// http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
	// https://splash.riverbed.com/docs/DOC-1926
	if ip := net.ParseIP(r.Header.Get("x-cluster-client-ip")); ip != nil {
		return ip
	}

	if ip := net.ParseIP(r.Header.Get("x-forwarded")); ip != nil {
		return ip
	}

	if ip := net.ParseIP(r.Header.Get("forwarded-for")); ip != nil {
		return ip
	}

	remoteAddr := r.RemoteAddr
	if raddr, ok := splitHostPort(remoteAddr); ok {
		remoteAddr = raddr
	}

	return net.ParseIP(remoteAddr)
}

func fromXForwardedFor(xfwdfor string) net.IP {
	// x-forwarded-for may return multiple IP addresses in the format:
	// "client IP, proxy 1 IP, proxy 2 IP"
	// Therefore, the right-most IP address is the IP address of the most recent proxy
	// and the left-most IP address is the IP address of the originating client.
	// source: http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/x-forwarded-headers.html
	// Azure Web App's also adds a port for some reason, so we'll only use the first part (the IP)
	for _, ip := range strings.Split(xfwdfor, ",") {
		ip = strings.TrimSpace(ip)
		if raddr, ok := splitHostPort(ip); ok {
			ip = raddr
		}

		// Sometimes IP addresses in this header can be 'unknown' (http://stackoverflow.com/a/11285650).
		// Therefore taking the left-most IP address that is not unknown
		// A Squid configuration directive can also set the value to "unknown" (http://www.squid-cache.org/Doc/config/forwarded_for/)
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return parsedIP
		}
	}

	return nil
}

func splitHostPort(addr string) (string, bool) {
	raddr, _, err := net.SplitHostPort(addr)
	return raddr, raddr != "" && err == nil
}
